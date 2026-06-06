package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"todotxt/internal/task"
)

type Store struct {
	TodoFile string
	DoneFile string
}

func New(todoFile string) *Store {
	if todoFile == "" {
		if dir := os.Getenv("TODO_DIR"); dir != "" {
			todoFile = filepath.Join(dir, "todo.txt")
		} else {
			todoFile = "todo.txt"
		}
	}
	doneFile := strings.TrimSuffix(todoFile, ".txt") + ".done.txt"
	if todoFile == "todo.txt" {
		doneFile = "done.txt"
	}
	return &Store{TodoFile: todoFile, DoneFile: doneFile}
}

func (s *Store) Load() ([]*task.Task, error) {
	file, err := os.Open(s.TodoFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []*task.Task{}, nil
		}
		return nil, fmt.Errorf("erro ao abrir %s: %w", s.TodoFile, err)
	}
	defer file.Close()

	var tasks []*task.Task
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimRight(scanner.Text(), "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		tasks = append(tasks, task.Parse(line, lineNum))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler %s: %w", s.TodoFile, err)
	}
	return tasks, nil
}

func (s *Store) Save(tasks []*task.Task) error {
	lines := make([]string, 0, len(tasks))
	for _, t := range tasks {
		if strings.TrimSpace(t.Description) == "" {
			continue
		}
		lines = append(lines, t.Format())
	}
	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n"
	}
	return os.WriteFile(s.TodoFile, []byte(content), 0644)
}

func (s *Store) Archive() (int, error) {
	tasks, err := s.Load()
	if err != nil {
		return 0, err
	}

	var pending []*task.Task
	var completed []*task.Task
	for _, t := range tasks {
		if t.Completed {
			completed = append(completed, t)
		} else {
			pending = append(pending, t)
		}
	}

	if len(completed) == 0 {
		return 0, nil
	}

	archiveFile, err := os.OpenFile(s.DoneFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("erro ao abrir %s: %w", s.DoneFile, err)
	}
	defer archiveFile.Close()

	writer := bufio.NewWriter(archiveFile)
	for _, t := range completed {
		if _, err := writer.WriteString(t.Format() + "\n"); err != nil {
			return 0, err
		}
	}
	if err := writer.Flush(); err != nil {
		return 0, err
	}

	if err := s.Save(pending); err != nil {
		return 0, err
	}
	return len(completed), nil
}

func (s *Store) Exists() bool {
	_, err := os.Stat(s.TodoFile)
	return err == nil
}
