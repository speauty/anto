package core

import (
	"anto/common"
	"anto/domain/repository"
	"anto/domain/service/translator"
	"anto/lib/srt"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSubscribeStartsWithTaskSnapshot(t *testing.T) {
	service := NewService()
	events, unsubscribe := service.Subscribe()
	defer unsubscribe()

	event := <-events
	if event.Type != "tasks.snapshot" {
		t.Fatalf("expected tasks.snapshot, got %q", event.Type)
	}
	if event.Tasks == nil {
		t.Fatal("expected a non-nil task snapshot")
	}
	if len(event.Tasks) != 0 {
		t.Fatalf("expected empty task snapshot, got %d tasks", len(event.Tasks))
	}
}

func TestCloneTaskCopiesOutputFiles(t *testing.T) {
	original := &Task{ID: "task-1", OutputFiles: []string{"one.srt"}}
	copy := cloneTask(original)
	copy.OutputFiles[0] = "changed.srt"
	if original.OutputFiles[0] != "one.srt" {
		t.Fatal("cloneTask shares the output files backing array")
	}
}

func TestTaskOutputMetadataIsRecorded(t *testing.T) {
	service := NewService()
	service.tasks["task-1"] = &Task{ID: "task-1"}
	service.setOutputDirectory("task-1", "C:/anto")
	service.addOutputFile("task-1", "C:/anto/source.srt")
	service.addOutputFile("task-1", "C:/anto/source.srt")
	task, ok := service.GetTask("task-1")
	if !ok || task.OutputDirectory != "C:/anto" || len(task.OutputFiles) != 1 {
		t.Fatalf("unexpected output metadata: %+v", task)
	}
}

func TestUpdatePublishesTranslatingEvent(t *testing.T) {
	service := NewService()
	service.tasks["task-1"] = &Task{ID: "task-1"}
	events, unsubscribe := service.Subscribe()
	defer unsubscribe()
	<-events // initial snapshot
	service.update("task-1", TaskTranslating, "正在翻译字幕（1/1）", "source.srt")
	event := <-events
	if event.Type != "task.translating" || event.Task == nil || event.Task.Status != TaskTranslating {
		t.Fatalf("unexpected translating event: %+v", event)
	}
}

func TestAppConfigDefaultsAreStable(t *testing.T) {
	app := appConfigFromSettings(map[string]interface{}{})
	if app.TrackExport != 2 || !app.KeepOriginal || !app.NotifyOnComplete || !app.AutoConnect {
		t.Fatalf("unexpected defaults: %+v", app)
	}
	if app.Concurrency != 5 || app.LogLevel != "信息" {
		t.Fatalf("unexpected numeric defaults: %+v", app)
	}
}

func TestProviderSecretsAreRedacted(t *testing.T) {
	values := redactValues(map[string]interface{}{"api_key": "secret", "region": "cn"})
	if values["api_key"] != redactedSecret || values["region"] != "cn" {
		t.Fatalf("unexpected redaction: %+v", values)
	}
}

type batchServiceTranslator struct{ texts []string }

func (b *batchServiceTranslator) TranslateBatch(_ context.Context, _, _ string, texts []string) ([]string, error) {
	b.texts = append([]string(nil), texts...)
	return []string{"第一个", "第二个"}, nil
}

func TestBatchSubtitleResultsUseBlockPositionAndRespectDeltaMode(t *testing.T) {
	service := NewService()
	service.tasks["batch-test"] = &Task{ID: "batch-test"}
	subtitles := &srt.Srt{Blocks: []*srt.Block{
		{MainTrack: "same"},
		{MainTrack: "skip", SubTrack: "已有译文"},
		{MainTrack: "same"},
	}}
	engine := &batchServiceTranslator{}
	err := service.translateBatchBlocks(context.Background(), "batch-test", engine,
		TaskRequest{FromLang: "en", ToLang: "zh", TranslateMode: common.ModeDelta}, "source.srt", subtitles)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(engine.texts, []string{"same", "same"}) ||
		subtitles.Blocks[0].SubTrack != "第一个" ||
		subtitles.Blocks[1].SubTrack != "已有译文" ||
		subtitles.Blocks[2].SubTrack != "第二个" {
		t.Fatalf("batch subtitles were misaligned: texts=%q blocks=%+v", engine.texts, subtitles.Blocks)
	}
}

func TestCreateTaskRejectsUnavailableTranslator(t *testing.T) {
	repository.GetTranslators().Register(&serviceTestTranslator{
		id:    "unavailable-test",
		valid: false,
	})
	service := NewService()
	file := writeServiceTestSRT(t)
	_, err := service.CreateTask(TaskRequest{
		Files:        []string{file},
		TranslatorID: "unavailable-test",
		FromLang:     "en",
		ToLang:       "zh",
	})
	if err == nil || err.Error() != "翻译引擎未配置或不可用: unavailable-test" {
		t.Fatalf("unexpected unavailable translator error: %v", err)
	}
}

func TestCreateTaskRejectsSameLanguagePair(t *testing.T) {
	repository.GetTranslators().Register(&serviceTestTranslator{
		id:    "valid-test",
		valid: true,
	})
	service := NewService()
	file := writeServiceTestSRT(t)
	_, err := service.CreateTask(TaskRequest{
		Files:        []string{file},
		TranslatorID: "valid-test",
		FromLang:     "EN",
		ToLang:       "en",
	})
	if err == nil || err.Error() != "源语言和目标语言不能相同" {
		t.Fatalf("unexpected same-language error: %v", err)
	}
}

func TestCreateTaskDeltaModePreservesExistingSubTrack(t *testing.T) {
	repository.GetTranslators().Register(&serviceTestTranslator{
		id:    "delta-test",
		valid: true,
	})
	input := filepath.Join(t.TempDir(), "source.srt")
	content := "1\n00:00:00,000 --> 00:00:01,000\nhello\n已有译文\n\n2\n00:00:01,000 --> 00:00:02,000\nworld\n"
	if err := os.WriteFile(input, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "output")
	service := NewService()
	task, err := service.CreateTask(TaskRequest{
		Files:           []string{input},
		TranslatorID:    "delta-test",
		FromLang:        "en",
		ToLang:          "zh",
		TranslateMode:   common.ModeDelta,
		FlagTrackExport: 2,
		OutputDirectory: output,
	})
	if err != nil {
		t.Fatal(err)
	}
	for attempt := 0; attempt < 50; attempt++ {
		current, ok := service.GetTask(task.ID)
		if !ok {
			t.Fatal("created task disappeared")
		}
		if current.Status == TaskCompleted {
			break
		}
		if current.Status == TaskFailed || current.Status == TaskCancelled {
			t.Fatalf("delta task ended unexpectedly: %+v", current)
		}
		time.Sleep(10 * time.Millisecond)
	}
	current, _ := service.GetTask(task.ID)
	if current.Status != TaskCompleted {
		t.Fatalf("delta task did not complete: %+v", current)
	}
	if len(current.OutputFiles) != 1 || current.OutputDirectory != output {
		t.Fatalf("unexpected output metadata: %+v", current)
	}
	encoded, err := os.ReadFile(current.OutputFiles[0])
	if err != nil {
		t.Fatal(err)
	}
	result := string(encoded)
	if !strings.Contains(result, "已有译文") || !strings.Contains(result, "译文:world") {
		t.Fatalf("delta output did not preserve and translate tracks: %q", result)
	}
	if strings.Contains(result, "\nhello\n") || strings.Contains(result, "\nworld\n") {
		t.Fatalf("副轨 output unexpectedly contains source track: %q", result)
	}
}

func TestCreateTaskFailsWhenOutputPathIsFile(t *testing.T) {
	repository.GetTranslators().Register(&serviceTestTranslator{
		id:    "output-path-test",
		valid: true,
	})
	input := writeServiceTestSRT(t)
	output := filepath.Join(t.TempDir(), "output")
	if err := os.WriteFile(output, []byte("not a directory"), 0600); err != nil {
		t.Fatal(err)
	}
	service := NewService()
	task, err := service.CreateTask(TaskRequest{
		Files:           []string{input},
		TranslatorID:    "output-path-test",
		FromLang:        "en",
		ToLang:          "zh",
		OutputDirectory: output,
	})
	if err != nil {
		t.Fatal(err)
	}
	for attempt := 0; attempt < 50; attempt++ {
		current, ok := service.GetTask(task.ID)
		if !ok {
			t.Fatal("created task disappeared")
		}
		if current.Status == TaskFailed {
			if !strings.Contains(current.Error, "创建输出目录失败") {
				t.Fatalf("unexpected output failure: %+v", current)
			}
			return
		}
		if current.Status == TaskCompleted || current.Status == TaskCancelled {
			t.Fatalf("task ended unexpectedly: %+v", current)
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("task did not reach failed status")
}

func TestCreateTaskExportsAllTrackModes(t *testing.T) {
	repository.GetTranslators().Register(&serviceTestTranslator{
		id:    "track-mode-test",
		valid: true,
	})
	for _, testCase := range []struct {
		name            string
		trackMode       int
		wantSource      bool
		wantTranslation bool
	}{
		{name: "all tracks", trackMode: 0, wantSource: true, wantTranslation: true},
		{name: "main track", trackMode: 1, wantSource: true},
		{name: "sub track", trackMode: 2, wantTranslation: true},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			input := writeServiceTestSRT(t)
			output := filepath.Join(t.TempDir(), "output")
			service := NewService()
			task, err := service.CreateTask(TaskRequest{
				Files:           []string{input},
				TranslatorID:    "track-mode-test",
				FromLang:        "en",
				ToLang:          "zh",
				FlagTrackExport: testCase.trackMode,
				OutputDirectory: output,
			})
			if err != nil {
				t.Fatal(err)
			}
			current := waitForServiceTask(t, service, task.ID)
			if current.Status != TaskCompleted {
				t.Fatalf("task did not complete: %+v", current)
			}
			encoded, err := os.ReadFile(current.OutputFiles[0])
			if err != nil {
				t.Fatal(err)
			}
			result := string(encoded)
			hasSource := strings.Contains(result, "\nhello\n")
			if hasSource != testCase.wantSource {
				t.Fatalf("source track presence = %v, output: %q", hasSource, result)
			}
			if strings.Contains(result, "译文:hello") != testCase.wantTranslation {
				t.Fatalf("translated track presence = %v, output: %q", strings.Contains(result, "译文:hello"), result)
			}
		})
	}
}

func waitForServiceTask(t *testing.T, service *Service, id string) *Task {
	t.Helper()
	for attempt := 0; attempt < 50; attempt++ {
		current, ok := service.GetTask(id)
		if !ok {
			t.Fatal("created task disappeared")
		}
		if current.Status == TaskCompleted || current.Status == TaskFailed || current.Status == TaskCancelled {
			return current
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("task did not reach a terminal status")
	return nil
}

func writeServiceTestSRT(t *testing.T) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "source.srt")
	if err := os.WriteFile(file, []byte("1\n00:00:00,000 --> 00:00:01,000\nhello\n"), 0600); err != nil {
		t.Fatal(err)
	}
	return file
}

type serviceTestTranslator struct {
	id    string
	valid bool
}

func (t *serviceTestTranslator) Init(translator.ImplConfig) {}
func (t *serviceTestTranslator) GetId() string              { return t.id }
func (t *serviceTestTranslator) GetShortId() string         { return t.id }
func (t *serviceTestTranslator) GetName() string            { return t.id }
func (t *serviceTestTranslator) GetCfg() translator.ImplConfig {
	return &translator.DefaultConfig{}
}
func (t *serviceTestTranslator) GetLangSupported() []translator.LangPair { return nil }
func (t *serviceTestTranslator) GetSep() string                          { return "" }
func (t *serviceTestTranslator) IsValid() bool                           { return t.valid }
func (t *serviceTestTranslator) Translate(_ context.Context, args *translator.TranslateArgs) (*translator.TranslateRes, error) {
	return &translator.TranslateRes{Results: []*translator.TranslateResBlock{{Id: args.TextContent, TextTranslated: "译文:" + args.TextContent}}}, nil
}
