package core

import (
	"anto/cfg"
	"anto/common"
	"anto/domain/repository"
	"anto/domain/service/translator"
	"anto/domain/service/translator/caiyunai"
	"anto/domain/service/translator/deepl_pro"
	"anto/domain/service/translator/youdao_batch"
	"anto/lib/srt"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Service struct {
	mu       sync.RWMutex
	tasks    map[string]*Task
	cancels  map[string]context.CancelFunc
	configMu sync.Mutex
	subsMu   sync.RWMutex
	subs     map[chan Event]struct{}
}

var ErrTranslatorNotFound = errors.New("翻译引擎不存在")
var ErrConnectionTestUnsupported = errors.New("翻译引擎不支持连接测试")

func NewService() *Service {
	return &Service{tasks: make(map[string]*Task), cancels: make(map[string]context.CancelFunc), subs: make(map[chan Event]struct{})}
}

// Close cancels active work before the transport is shut down.
func (s *Service) Close() {
	s.mu.RLock()
	cancels := make([]context.CancelFunc, 0, len(s.cancels))
	for _, cancel := range s.cancels {
		cancels = append(cancels, cancel)
	}
	s.mu.RUnlock()
	for _, cancel := range cancels {
		cancel()
	}
}

func (s *Service) Translators() []*TranslatorInfo {
	result := make([]*TranslatorInfo, 0)
	for _, item := range repository.GetTranslators().GetNamesAll() {
		engine := repository.GetTranslators().GetById(item.Key)
		info := &TranslatorInfo{ID: item.Key, Name: item.Name, Valid: engine != nil && engine.IsValid()}
		if engine != nil {
			seen := make(map[string]struct{})
			for _, language := range engine.GetLangSupported() {
				if language.Key == "" {
					continue
				}
				if _, exists := seen[language.Key]; exists {
					continue
				}
				seen[language.Key] = struct{}{}
				info.Languages = append(info.Languages, LanguageInfo{ID: language.Key, Name: language.Name})
			}
		}
		result = append(result, info)
	}
	return result
}

func (s *Service) TestTranslatorConnection(ctx context.Context, id, region string) (*translator.ConnectionResult, error) {
	engine := repository.GetTranslators().GetById(id)
	if engine == nil {
		return nil, fmt.Errorf("%w: %s", ErrTranslatorNotFound, id)
	}
	tester, ok := engine.(translator.ConnectionTester)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrConnectionTestUnsupported, id)
	}
	return tester.TestConnection(ctx, region)
}

func (s *Service) Config() ConfigDTO {
	s.configMu.Lock()
	defer s.configMu.Unlock()
	return buildConfigDTO()
}

func (s *Service) UpdateConfig(update ConfigUpdate) error {
	s.configMu.Lock()
	defer s.configMu.Unlock()
	viper := cfg.Singleton().GetViper()
	if len(update.App) != 0 && string(update.App) != "null" {
		app := appConfigFromSettings(viper.AllSettings())
		if err := json.Unmarshal(update.App, &app); err != nil {
			return &ConfigValidationError{Message: fmt.Sprintf("应用设置格式无效: %v", err)}
		}
		if err := validateAppConfig(app); err != nil {
			return err
		}
		viper.Set("ui", appConfigValues(app))
	}
	if update.UI != nil {
		app := appConfigFromSettings(map[string]interface{}{"ui": update.UI})
		if err := validateAppConfig(app); err != nil {
			return err
		}
		viper.Set("ui", update.UI)
	}
	for provider, values := range update.Providers {
		if repository.GetTranslators().GetById(provider) == nil {
			return &ConfigValidationError{Message: fmt.Sprintf("翻译引擎不存在: %s", provider)}
		}
		delete(values, "configured")
		delete(values, "values")
		if err := validateProviderValues(provider, values); err != nil {
			return err
		}
		// A redacted secret means the client left the current value untouched.
		current, _ := viper.Get(provider).(map[string]interface{})
		for key, value := range values {
			if stringValue(value) == redactedSecret {
				delete(values, key)
			}
		}
		for key, value := range current {
			if key == "configured" || key == "values" {
				continue
			}
			if _, exists := values[key]; !exists {
				values[key] = value
			}
		}
		viper.Set(provider, values)
	}
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	return nil
}

type TranslatorInfo struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Valid     bool           `json:"valid"`
	Languages []LanguageInfo `json:"languages,omitempty"`
}

type LanguageInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Service) GetTask(id string) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	if !ok {
		return nil, false
	}
	copy := *task
	return &copy, true
}

func (s *Service) ListTasks() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, cloneTask(task))
	}
	return result
}

func (s *Service) Subscribe() (<-chan Event, func()) {
	channel := make(chan Event, 64)
	// Capture the snapshot before registering the subscriber. Registration and
	// enqueueing the snapshot happen under the same lock, so a concurrent
	// publisher cannot put an incremental event ahead of the initial snapshot.
	s.mu.RLock()
	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, cloneTask(task))
	}
	s.mu.RUnlock()
	s.subsMu.Lock()
	s.subs[channel] = struct{}{}
	s.enqueue(channel, Event{Type: "tasks.snapshot", Tasks: tasks})
	s.subsMu.Unlock()
	return channel, func() {
		s.subsMu.Lock()
		if _, ok := s.subs[channel]; ok {
			delete(s.subs, channel)
			close(channel)
		}
		s.subsMu.Unlock()
	}
}

func (s *Service) CreateTask(request TaskRequest) (*Task, error) {
	files, err := collectFiles(request)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("未找到需要翻译的 SRT 文件")
	}
	engine := repository.GetTranslators().GetById(request.TranslatorID)
	if engine == nil {
		return nil, fmt.Errorf("翻译引擎不存在: %s", request.TranslatorID)
	}
	if !engine.IsValid() {
		return nil, fmt.Errorf("翻译引擎未配置或不可用: %s", request.TranslatorID)
	}
	if strings.TrimSpace(request.FromLang) != "" &&
		strings.EqualFold(strings.TrimSpace(request.FromLang), strings.TrimSpace(request.ToLang)) {
		return nil, errors.New("源语言和目标语言不能相同")
	}
	if request.TranslatorID == "caiyunai" && !caiyunai.SupportsDirection(request.FromLang, request.ToLang) {
		return nil, fmt.Errorf("彩云小译不支持 %s 到 %s 的翻译方向", request.FromLang, request.ToLang)
	}
	if request.TranslatorID == "deepl_pro" && !deepl_pro.SupportsDirection(request.FromLang, request.ToLang) {
		return nil, fmt.Errorf("DeepL Pro 不支持 %s 到 %s 的翻译方向", request.FromLang, request.ToLang)
	}
	if request.TranslatorID == "youdao_batch" && !youdao_batch.SupportsDirection(request.FromLang, request.ToLang) {
		return nil, fmt.Errorf("有道智云不支持 %s 到 %s 的翻译方向", request.FromLang, request.ToLang)
	}
	if request.TranslateMode == "" {
		request.TranslateMode = common.ModeFull
	}
	if request.MainTrackReport == "" {
		request.MainTrackReport = common.LangDirectionFrom
	}
	id := newID()
	now := time.Now()
	task := &Task{
		ID:              id,
		Status:          TaskQueued,
		Total:           len(files),
		OutputDirectory: request.OutputDirectory,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	s.mu.Lock()
	s.tasks[id] = task
	created := cloneTask(task)
	s.mu.Unlock()
	s.publish(Event{Type: "task.created", TaskID: id, Task: created})
	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.cancels[id] = cancel
	s.mu.Unlock()
	go s.run(ctx, id, request, files)
	return created, nil
}

func (s *Service) CancelTask(id string) error {
	s.mu.RLock()
	cancel, ok := s.cancels[id]
	s.mu.RUnlock()
	if !ok {
		return errors.New("任务不存在或已结束")
	}
	cancel()
	return nil
}

func (s *Service) run(ctx context.Context, id string, request TaskRequest, files []string) {
	defer func() { s.mu.Lock(); delete(s.cancels, id); s.mu.Unlock() }()
	engine := repository.GetTranslators().GetById(request.TranslatorID)
	for index, file := range files {
		if ctx.Err() != nil {
			s.update(id, TaskCancelled, "任务已取消", "")
			return
		}
		s.update(id, TaskReading, "正在读取字幕", file)
		if err := s.translateFile(ctx, id, engine, request, file); err != nil {
			if ctx.Err() != nil {
				s.update(id, TaskCancelled, "任务已取消", file)
			} else {
				s.update(id, TaskFailed, err.Error(), file)
			}
			return
		}
		s.mu.Lock()
		if task := s.tasks[id]; task != nil {
			task.Completed = index + 1
			task.UpdatedAt = time.Now()
		}
		snapshot := cloneTask(s.tasks[id])
		s.mu.Unlock()
		s.publish(Event{Type: "task.progress", TaskID: id, Task: snapshot})
	}
	s.update(id, TaskCompleted, "翻译完成", "")
}

func (s *Service) translateFile(ctx context.Context, taskID string, engine translator.ImplTranslator, request TaskRequest, file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}
	srtFile := new(srt.Srt)
	srtFile.FilePath = file
	srtFile.FileSize = len(data)
	srtFile.FileNameSync()
	if err = srtFile.Decode(bytes.NewReader(data)); err != nil {
		return fmt.Errorf("解析文件失败: %w", err)
	}
	srtFile.CntBlock = len(srtFile.Blocks)
	if batch, ok := engine.(translator.BatchTranslator); ok {
		if err := s.translateBatchBlocks(ctx, taskID, batch, request, file, srtFile); err != nil {
			return err
		}
	} else {
		if err := s.translateBlocks(ctx, taskID, engine, request, file, srtFile); err != nil {
			return err
		}
	}
	output := request.OutputDirectory
	if output == "" {
		output = filepath.Join(filepath.Dir(strings.TrimSuffix(file, filepath.Ext(file))), common.AppName)
	}
	s.setOutputDirectory(taskID, output)
	if err = os.MkdirAll(output, os.ModePerm); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}
	encoded, err := srtFile.Encode(&srt.EncodeOpt{FlagIsInverse: request.MainTrackReport == common.LangDirectionTo, FlagTrackExport: request.FlagTrackExport})
	if err != nil {
		return fmt.Errorf("编码字幕失败: %w", err)
	}
	outputFile := filepath.Join(output, strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))+filepath.Ext(file))
	s.update(taskID, TaskWriting, "正在写入字幕", outputFile)
	if err = os.WriteFile(outputFile, encoded, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	s.addOutputFile(taskID, outputFile)
	return nil
}

func (s *Service) translateBatchBlocks(ctx context.Context, taskID string, engine translator.BatchTranslator, request TaskRequest, file string, srtFile *srt.Srt) error {
	indexes := make([]int, 0, len(srtFile.Blocks))
	texts := make([]string, 0, len(srtFile.Blocks))
	for index, block := range srtFile.Blocks {
		if block.SubTrack != "" && request.TranslateMode == common.ModeDelta {
			continue
		}
		indexes = append(indexes, index)
		texts = append(texts, block.MainTrack)
	}
	if len(texts) == 0 {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	s.update(taskID, TaskTranslating, fmt.Sprintf("正在批量翻译 %d 条字幕", len(texts)), file)
	translated, err := engine.TranslateBatch(ctx, request.FromLang, request.ToLang, texts)
	if err != nil {
		return fmt.Errorf("翻译失败: %w", err)
	}
	if len(translated) != len(texts) {
		return fmt.Errorf("翻译失败: 返回 %d 条字幕，预期 %d 条", len(translated), len(texts))
	}
	for index, value := range translated {
		srtFile.Blocks[indexes[index]].SubTrack = value
	}
	return nil
}

func (s *Service) translateBlocks(ctx context.Context, taskID string, engine translator.ImplTranslator, request TaskRequest, file string, srtFile *srt.Srt) error {
	chunks := make([]string, 0)
	current := ""
	max := engine.GetCfg().GetMaxCharNum()
	appendChunk := func(value string) {
		if value != "" {
			chunks = append(chunks, value)
		}
	}
	for _, block := range srtFile.Blocks {
		if block.SubTrack != "" && request.TranslateMode == common.ModeDelta {
			continue
		}
		if engine.GetSep() == "" {
			chunks = append(chunks, block.MainTrack)
			continue
		}
		if max > 0 && len(current) >= max {
			appendChunk(current)
			current = ""
		}
		if max > 0 && len(block.MainTrack) >= max {
			appendChunk(current)
			current = ""
			chunks = append(chunks, block.MainTrack)
			continue
		}
		separator := engine.GetSep()
		if max > 0 && len(current)+len(block.MainTrack)+len(separator) > max {
			appendChunk(current)
			current = ""
		}
		if current == "" {
			current = block.MainTrack
		} else {
			current += separator + block.MainTrack
		}
	}
	appendChunk(current)
	if len(chunks) == 0 {
		return nil
	}
	for index, chunk := range chunks {
		if err := ctx.Err(); err != nil {
			return err
		}
		s.update(taskID, TaskTranslating,
			fmt.Sprintf("正在翻译字幕（%d/%d）", index+1, len(chunks)), file)
		result, err := engine.Translate(ctx, &translator.TranslateArgs{FromLang: request.FromLang, ToLang: request.ToLang, TextContent: chunk})
		if err != nil {
			return fmt.Errorf("翻译失败: %w", err)
		}
		for _, translated := range result.Results {
			for _, block := range srtFile.Blocks {
				if block.MainTrack == translated.Id {
					block.SubTrack = translated.TextTranslated
				}
			}
		}
	}
	return nil
}

func (s *Service) setOutputDirectory(id, directory string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task := s.tasks[id]; task != nil && task.OutputDirectory == "" {
		task.OutputDirectory = directory
		task.UpdatedAt = time.Now()
	}
}

func (s *Service) addOutputFile(id, file string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task := s.tasks[id]
	if task == nil {
		return
	}
	for _, existing := range task.OutputFiles {
		if existing == file {
			return
		}
	}
	task.OutputFiles = append(task.OutputFiles, file)
	task.UpdatedAt = time.Now()
}

func (s *Service) update(id string, status TaskStatus, message, file string) {
	s.mu.Lock()
	task := s.tasks[id]
	if task == nil {
		s.mu.Unlock()
		return
	}
	task.Status = status
	task.CurrentFile = file
	task.UpdatedAt = time.Now()
	task.Error = ""
	if status == TaskFailed {
		task.Error = message
	}
	snapshot := cloneTask(task)
	s.mu.Unlock()
	s.publish(Event{Type: "task." + string(status), TaskID: id, Message: message, Task: snapshot})
}

func (s *Service) publish(event Event) {
	s.subsMu.RLock()
	defer s.subsMu.RUnlock()
	for channel := range s.subs {
		s.enqueue(channel, event)
	}
}

func (s *Service) enqueue(channel chan Event, event Event) {
	select {
	case channel <- event:
	default:
		// Keep the stream live when a client is slow. The next reconnect also
		// receives a full snapshot, so dropping the oldest buffered event is
		// safer than blocking every active task.
		select {
		case <-channel:
		default:
		}
		select {
		case channel <- event:
		default:
		}
	}
}

func cloneTask(task *Task) *Task {
	if task == nil {
		return nil
	}
	copy := *task
	copy.OutputFiles = append([]string(nil), task.OutputFiles...)
	return &copy
}

func collectFiles(request TaskRequest) ([]string, error) {
	files := make([]string, 0, len(request.Files))
	seen := make(map[string]struct{})
	add := func(path string) {
		if path == "" || strings.ToLower(filepath.Ext(path)) != ".srt" {
			return
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return
		}
		if _, ok := seen[abs]; !ok {
			seen[abs] = struct{}{}
			files = append(files, abs)
		}
	}
	for _, file := range request.Files {
		add(file)
	}
	if request.Directory != "" {
		entries, err := os.ReadDir(request.Directory)
		if err != nil {
			return nil, fmt.Errorf("读取目录失败: %w", err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				add(filepath.Join(request.Directory, entry.Name()))
			}
		}
	}
	return files, nil
}

func newID() string {
	data := make([]byte, 12)
	if _, err := rand.Read(data); err != nil {
		return fmt.Sprintf("task-%d", time.Now().UnixNano())
	}
	return "task-" + hex.EncodeToString(data)
}
