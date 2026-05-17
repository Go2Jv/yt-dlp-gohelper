package history

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Status string

const (
	StatusRunning Status = "running"
	StatusFailed  Status = "failed"
)

type Task struct {
	ID            string    `json:"id"`
	URL           string    `json:"url"`
	Args          []string  `json:"args"`
	OutputDir     string    `json:"outputDir"`
	OutputTmpl    string    `json:"outputTemplate"`
	Status        Status    `json:"status"`
	Error         string    `json:"error,omitempty"`
	FragmentsTotal int      `json:"fragmentsTotal,omitempty"`
	TempFiles     []string  `json:"tempFiles,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Store struct {
	path string
}

func NewStore() *Store {
	return &Store{path: defaultPath()}
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) Load() ([]Task, error) {
	b, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(b, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Store) Save(tasks []Task) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *Store) Upsert(task Task) error {
	tasks, err := s.Load()
	if err != nil {
		return err
	}
	now := time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = now
	}
	task.UpdatedAt = now

	for i := range tasks {
		if tasks[i].ID == task.ID {
			tasks[i] = task
			return s.Save(tasks)
		}
	}
	tasks = append(tasks, task)
	return s.Save(tasks)
}

func (s *Store) Delete(id string) error {
	tasks, err := s.Load()
	if err != nil {
		return err
	}
	out := tasks[:0]
	for i := range tasks {
		if tasks[i].ID != id {
			out = append(out, tasks[i])
		}
	}
	return s.Save(out)
}

func NewID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func defaultPath() string {
	dir, err := os.UserConfigDir()
	if err == nil && dir != "" {
		return filepath.Join(dir, "yt-dlp-gohelper", "history.json")
	}
	wd, _ := os.Getwd()
	return filepath.Join(wd, ".yt-dlp-gohelper.history.json")
}

