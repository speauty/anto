package core

import (
	"anto/common"
	"time"
)

// TaskRequest is the platform-neutral description of one translation batch.
// Files and Directory are both accepted; when Directory is set, only SRT files
// directly inside it are added to the task.
type TaskRequest struct {
	Files           []string             `json:"files"`
	Directory       string               `json:"directory"`
	TranslatorID    string               `json:"translator_id"`
	FromLang        string               `json:"from_lang"`
	ToLang          string               `json:"to_lang"`
	TranslateMode   common.TranslateMode `json:"translate_mode"`
	MainTrackReport common.LangDirection `json:"main_track_report"`
	FlagTrackExport int                  `json:"flag_track_export"`
	OutputDirectory string               `json:"output_directory"`
}

type TaskStatus string

const (
	TaskQueued      TaskStatus = "queued"
	TaskReading     TaskStatus = "reading"
	TaskTranslating TaskStatus = "translating"
	TaskWriting     TaskStatus = "writing"
	TaskCompleted   TaskStatus = "completed"
	TaskFailed      TaskStatus = "failed"
	TaskCancelled   TaskStatus = "cancelled"
)

type Task struct {
	ID              string     `json:"id"`
	Status          TaskStatus `json:"status"`
	Total           int        `json:"total"`
	Completed       int        `json:"completed"`
	CurrentFile     string     `json:"current_file,omitempty"`
	OutputDirectory string     `json:"output_directory,omitempty"`
	OutputFiles     []string   `json:"output_files,omitempty"`
	Error           string     `json:"error,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Event struct {
	Type    string  `json:"type"`
	TaskID  string  `json:"task_id,omitempty"`
	Message string  `json:"message,omitempty"`
	Task    *Task   `json:"task,omitempty"`
	Tasks   []*Task `json:"tasks,omitempty"`
}
