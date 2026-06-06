package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"todotxt/internal/task"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}
	return string(b)
}

func TestLoadNonExistent(t *testing.T) {
	dir := t.TempDir()
	s := New(filepath.Join(dir, "todo.txt"))
	tasks, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestLoadExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todo.txt")
	writeFile(t, path, "(A) 2024-01-15 Call Mom +Family @phone\nx 2024-01-10 2024-01-08 Write report +Work\n")

	s := New(path)
	tasks, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Priority != "A" {
		t.Errorf("task[0].Priority: got %q, want A", tasks[0].Priority)
	}
	if !tasks[1].Completed {
		t.Error("task[1] should be completed")
	}
}

func TestLoadSkipsEmptyLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todo.txt")
	writeFile(t, path, "Task 1\n\n\nTask 2\n\n")

	s := New(path)
	tasks, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestSaveAndReload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todo.txt")
	s := New(path)

	t1 := task.New("Comprar leite +Pessoal @supermercado")
	t1.SetPriority("A")
	t2 := task.New("Estudar Go +Estudo")
	t3 := task.New("Tarefa concluída")
	t3.MarkComplete()

	if err := s.Save([]*task.Task{t1, t2, t3}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	content := readFile(t, path)
	for _, want := range []string{"Comprar leite", "Estudar Go", "x "} {
		if !strings.Contains(content, want) {
			t.Errorf("expected %q in file, got:\n%s", want, content)
		}
	}

	tasks, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tasks))
	}
	if tasks[0].Priority != "A" {
		t.Errorf("reloaded priority: got %q", tasks[0].Priority)
	}
	if !tasks[2].Completed {
		t.Error("reloaded completion lost")
	}
}

func TestSaveEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todo.txt")
	s := New(path)

	if err := s.Save([]*task.Task{}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file should exist: %v", err)
	}
}

func TestArchive(t *testing.T) {
	dir := t.TempDir()
	todoPath := filepath.Join(dir, "todo.txt")
	donePath := filepath.Join(dir, "done.txt")
	s := New(todoPath)

	pending := task.New("Tarefa pendente")
	completed := task.New("Tarefa concluída")
	completed.MarkComplete()
	other := task.New("Outra pendente")

	if err := s.Save([]*task.Task{pending, completed, other}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	count, err := s.Archive()
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 archived, got %d", count)
	}

	pendingAfter, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(pendingAfter) != 2 {
		t.Errorf("expected 2 pending after archive, got %d", len(pendingAfter))
	}
	for _, tk := range pendingAfter {
		if tk.Completed {
			t.Error("completed task should not remain in todo.txt")
		}
	}

	doneContent := readFile(t, donePath)
	if !strings.Contains(doneContent, "Tarefa concluída") {
		t.Errorf("done.txt should contain completed task, got:\n%s", doneContent)
	}
}

func TestArchiveNoCompleted(t *testing.T) {
	dir := t.TempDir()
	todoPath := filepath.Join(dir, "todo.txt")
	donePath := filepath.Join(dir, "done.txt")
	s := New(todoPath)

	pending := task.New("Pendente")
	if err := s.Save([]*task.Task{pending}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	count, err := s.Archive()
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 archived, got %d", count)
	}
	if _, err := os.Stat(donePath); err == nil {
		t.Error("done.txt should not be created when no completed tasks")
	}
}

func TestArchiveAppends(t *testing.T) {
	dir := t.TempDir()
	todoPath := filepath.Join(dir, "todo.txt")
	donePath := filepath.Join(dir, "done.txt")
	s := New(todoPath)

	writeFile(t, donePath, "x 2024-01-01 2023-12-30 Tarefa antiga\n")

	t1 := task.New("Tarefa 1")
	t1.MarkComplete()
	t2 := task.New("Tarefa 2")
	t2.MarkComplete()
	if err := s.Save([]*task.Task{t1, t2}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := s.Archive(); err != nil {
		t.Fatalf("Archive: %v", err)
	}

	content := readFile(t, donePath)
	if !strings.Contains(content, "Tarefa antiga") {
		t.Error("existing done.txt content should be preserved")
	}
	if !strings.Contains(content, "Tarefa 1") || !strings.Contains(content, "Tarefa 2") {
		t.Errorf("new tasks should be appended, got:\n%s", content)
	}
}

func TestCustomPaths(t *testing.T) {
	dir := t.TempDir()
	custom := filepath.Join(dir, "my-tasks.txt")
	writeFile(t, custom, "Task A\nTask B\n")

	s := New(custom)
	tasks, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}

	if !strings.HasSuffix(s.DoneFile, ".done.txt") {
		t.Errorf("DoneFile should be derived from todo path, got %s", s.DoneFile)
	}
}

func TestDoneFileDefaultForTodoDotTxt(t *testing.T) {
	s := New("todo.txt")
	if s.DoneFile != "done.txt" {
		t.Errorf("DoneFile: got %q, want done.txt", s.DoneFile)
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "todo.txt")
	s := New(path)
	if s.Exists() {
		t.Error("Exists should be false for new file")
	}
	writeFile(t, path, "x")
	if !s.Exists() {
		t.Error("Exists should be true after writing file")
	}
}
