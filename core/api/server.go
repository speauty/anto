package api

import (
	"anto/core"
	"anto/domain/service/translator"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	core     *core.Service
	token    string
	server   *http.Server
	listener net.Listener
}

type Ready struct {
	Port  int    `json:"port"`
	Token string `json:"token"`
}

type APIError struct {
	Error APIErrorBody `json:"error"`
}

type APIErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func New(service *core.Service) (*Server, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("监听本地端口失败: %w", err)
	}
	token, err := randomToken()
	if err != nil {
		_ = listener.Close()
		return nil, err
	}
	server := &Server{core: service, token: token, listener: listener}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", server.health)
	mux.HandleFunc("/api/v1/translators", server.auth(server.translators))
	mux.HandleFunc("/api/v1/translators/", server.auth(server.translator))
	mux.HandleFunc("/api/v1/config", server.auth(server.config))
	mux.HandleFunc("/api/v1/tasks", server.auth(server.tasks))
	mux.HandleFunc("/api/v1/tasks/", server.auth(server.task))
	mux.HandleFunc("/api/v1/events", server.auth(server.events))
	mux.HandleFunc("/api/v1/shutdown", server.auth(server.shutdown))
	server.server = &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	return server, nil
}

func (s *Server) Ready() Ready {
	return Ready{Port: s.listener.Addr().(*net.TCPAddr).Port, Token: s.token}
}

func (s *Server) Serve() error {
	err := s.server.Serve(s.listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Close(ctx context.Context) error {
	s.core.Close()
	return s.server.Shutdown(ctx)
}

func (s *Server) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Anto-Token") != s.token {
			writeError(writer, http.StatusUnauthorized, "unauthorized", "X-Anto-Token 无效")
			return
		}
		next(writer, request)
	}
}

func (s *Server) health(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "请求方法不支持")
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) translators(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "请求方法不支持")
		return
	}
	writeJSON(writer, http.StatusOK, s.core.Translators())
}

func (s *Server) translator(writer http.ResponseWriter, request *http.Request) {
	parts := strings.Split(strings.Trim(strings.TrimPrefix(request.URL.Path, "/api/v1/translators/"), "/"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] != "test-connection" {
		writeError(writer, http.StatusNotFound, "not_found", "接口不存在")
		return
	}
	if request.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "请求方法不支持")
		return
	}
	var input struct {
		Region string `json:"region"`
	}
	if request.ContentLength != 0 {
		if err := decodeJSON(request, &input); err != nil {
			writeError(writer, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}
	result, err := s.core.TestTranslatorConnection(request.Context(), parts[0], input.Region)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrTranslatorNotFound):
			writeError(writer, http.StatusNotFound, "translator_not_found", err.Error())
		case errors.Is(err, core.ErrConnectionTestUnsupported):
			writeError(writer, http.StatusNotImplemented, "connection_test_unsupported", err.Error())
		case errors.Is(err, translator.ErrInvalidConnectionRegion):
			writeError(writer, http.StatusBadRequest, "invalid_region", err.Error())
		default:
			writeError(writer, http.StatusBadGateway, "translator_connection_failed", err.Error())
		}
		return
	}
	writeJSON(writer, http.StatusOK, result)
}

func (s *Server) config(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		writeJSON(writer, http.StatusOK, s.core.Config())
	case http.MethodPut:
		var values core.ConfigUpdate
		if err := decodeJSON(request, &values); err != nil {
			writeError(writer, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		if err := s.core.UpdateConfig(values); err != nil {
			var validationErr *core.ConfigValidationError
			if errors.As(err, &validationErr) {
				writeError(writer, http.StatusBadRequest, "invalid_config", err.Error())
				return
			}
			writeError(writer, http.StatusInternalServerError, "config_write_failed", err.Error())
			return
		}
		writeJSON(writer, http.StatusOK, s.core.Config())
	default:
		writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "请求方法不支持")
	}
}

func (s *Server) tasks(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		writeJSON(writer, http.StatusOK, s.core.ListTasks())
		return
	}
	if request.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "请求方法不支持")
		return
	}
	var payload core.TaskRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	task, err := s.core.CreateTask(payload)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "task_create_failed", err.Error())
		return
	}
	writeJSON(writer, http.StatusAccepted, task)
}

func (s *Server) task(writer http.ResponseWriter, request *http.Request) {
	id := TaskPathID(request.URL.Path)
	if id == "" {
		writeError(writer, http.StatusBadRequest, "task_id_required", "任务 ID 不能为空")
		return
	}
	if request.Method == http.MethodGet {
		task, ok := s.core.GetTask(id)
		if !ok {
			writeError(writer, http.StatusNotFound, "task_not_found", "任务不存在")
			return
		}
		writeJSON(writer, http.StatusOK, task)
		return
	}
	if request.Method == http.MethodPost && strings.HasSuffix(strings.TrimSuffix(request.URL.Path, "/"), "/cancel") {
		id = strings.TrimSuffix(id, "/cancel")
		if err := s.core.CancelTask(id); err != nil {
			writeError(writer, http.StatusBadRequest, "task_cancel_failed", err.Error())
			return
		}
		writeJSON(writer, http.StatusAccepted, map[string]string{"status": "cancelling"})
		return
	}
	writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "请求方法不支持")
}

func (s *Server) events(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "请求方法不支持")
		return
	}
	flusher, ok := writer.(http.Flusher)
	if !ok {
		writeError(writer, http.StatusInternalServerError, "streaming_unsupported", "当前服务器不支持 SSE")
		return
	}
	channel, unsubscribe := s.core.Subscribe()
	defer unsubscribe()
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	flusher.Flush()
	for {
		select {
		case <-request.Context().Done():
			return
		case event, ok := <-channel:
			if !ok {
				return
			}
			data, _ := json.Marshal(event)
			_, _ = fmt.Fprintf(writer, "event: %s\ndata: %s\n\n", event.Type, data)
			flusher.Flush()
		}
	}
}

func (s *Server) shutdown(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "请求方法不支持")
		return
	}
	writeJSON(writer, http.StatusAccepted, map[string]string{"status": "shutting_down"})
	go func() { _ = s.Close(context.Background()) }()
}

func decodeJSON(request *http.Request, target interface{}) error {
	defer request.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(request.Body, 2<<20))
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("请求 JSON 无效: %w", err)
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, status int, value interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func writeError(writer http.ResponseWriter, status int, code, message string) {
	writeJSON(writer, status, APIError{Error: APIErrorBody{Code: code, Message: message}})
}

func randomToken() (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func TaskPathID(path string) string {
	return strings.Trim(strings.TrimPrefix(path, "/api/v1/tasks/"), "/")
}
