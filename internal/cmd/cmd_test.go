package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"todotxt/internal/store"
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	t.Setenv("NO_COLOR", "1")
	dir := t.TempDir()
	return store.New(filepath.Join(dir, "todo.txt"))
}

func readStoreFile(t *testing.T, s *store.Store) string {
	t.Helper()
	b, err := os.ReadFile(s.TodoFile)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	return string(b)
}

func TestAdd(t *testing.T) {
	s := newTestStore(t)

	out, err := Add(s, []string{"Comprar", "leite", "+Pessoal", "@supermercado"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if !strings.Contains(out, "Tarefa adicionada") {
		t.Errorf("output should confirm add, got %q", out)
	}
	if !strings.Contains(out, "Comprar leite +Pessoal @supermercado") {
		t.Errorf("output should include description, got %q", out)
	}

	content := readStoreFile(t, s)
	if !strings.Contains(content, "Comprar leite +Pessoal @supermercado") {
		t.Errorf("file should contain task, got %q", content)
	}
}

func TestAddWithPriorityAndDue(t *testing.T) {
	s := newTestStore(t)

	_, err := Add(s, []string{"Reunião", "+Trabalho", "pri:A", "due:2026-12-31"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	content := readStoreFile(t, s)
	if !strings.Contains(content, "(A)") {
		t.Errorf("expected priority (A), got %q", content)
	}
	if !strings.Contains(content, "due:2026-12-31") {
		t.Errorf("expected due date, got %q", content)
	}
}

func TestAddEmptyArgs(t *testing.T) {
	s := newTestStore(t)
	_, err := Add(s, nil)
	if err == nil {
		t.Error("Add with no args should return error")
	}
}

func TestAddInvalidDue(t *testing.T) {
	s := newTestStore(t)
	_, err := Add(s, []string{"Tarefa", "due:not-a-date"})
	if err == nil {
		t.Error("Add with invalid due should return error")
	}
	if !strings.Contains(err.Error(), "data inválida") {
		t.Errorf("error should mention invalid date, got %v", err)
	}
}

func TestAddMultiple(t *testing.T) {
	s := newTestStore(t)
	for _, args := range [][]string{
		{"Primeira"},
		{"Segunda"},
		{"Terceira"},
	} {
		if _, err := Add(s, args); err != nil {
			t.Fatalf("Add %v: %v", args, err)
		}
	}

	tasks, err := s.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestListEmpty(t *testing.T) {
	s := newTestStore(t)
	out, err := List(s, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if !strings.Contains(out, "Nenhuma tarefa") {
		t.Errorf("expected 'Nenhuma tarefa', got %q", out)
	}
}

func TestListWithTasks(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Tarefa 1"})
	_, _ = Add(s, []string{"Tarefa 2"})

	out, err := List(s, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if !strings.Contains(out, "Tarefa 1") || !strings.Contains(out, "Tarefa 2") {
		t.Errorf("expected both tasks, got %q", out)
	}
	if !strings.Contains(out, "Total: 2") {
		t.Errorf("expected total count, got %q", out)
	}
}

func TestListWithProjectFilter(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Tarefa 1", "+Trabalho"})
	_, _ = Add(s, []string{"Tarefa 2", "+Pessoal"})

	out, err := List(s, []string{"+Trabalho"})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if !strings.Contains(out, "Tarefa 1") {
		t.Errorf("expected Tarefa 1 in output, got %q", out)
	}
	if strings.Contains(out, "Tarefa 2") {
		t.Errorf("Tarefa 2 should be filtered out, got %q", out)
	}
}

func TestListAllShowsCompleted(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Pendente"})
	_, _ = Add(s, []string{"Concluída"})
	_, _ = Do(s, []string{"2"})

	out, err := List(s, []string{"all"})
	if err != nil {
		t.Fatalf("List all: %v", err)
	}
	if !strings.Contains(out, "Pendente") || !strings.Contains(out, "Concluída") {
		t.Errorf("expected both, got %q", out)
	}
}

func TestDo(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Comprar leite"})

	out, err := Do(s, []string{"1"})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	if !strings.Contains(out, "concluída") {
		t.Errorf("expected confirmation, got %q", out)
	}

	content := readStoreFile(t, s)
	if !strings.HasPrefix(content, "x ") {
		t.Errorf("task should be marked complete in file, got %q", content)
	}
}

func TestDoInvalidNumber(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Task"})

	_, err := Do(s, []string{"99"})
	if err == nil {
		t.Error("Do with invalid number should return error")
	}
}

func TestDoAlreadyDone(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Task"})
	_, _ = Do(s, []string{"1"})

	_, err := Do(s, []string{"1"})
	if err == nil {
		t.Error("Do on already completed task should return error")
	}
}

func TestDoNoArgs(t *testing.T) {
	s := newTestStore(t)
	_, err := Do(s, nil)
	if err == nil {
		t.Error("Do with no args should return error")
	}
}

func TestUndo(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Task"})
	_, _ = Do(s, []string{"1"})

	out, err := Undo(s, []string{"1"})
	if err != nil {
		t.Fatalf("Undo: %v", err)
	}
	if !strings.Contains(out, "reaberta") {
		t.Errorf("expected 'reaberta' in output, got %q", out)
	}

	tasks, _ := s.Load()
	if tasks[0].Completed {
		t.Error("task should be incomplete after Undo")
	}
}

func TestUndoNotCompleted(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Task"})
	_, err := Undo(s, []string{"1"})
	if err == nil {
		t.Error("Undo on pending task should return error")
	}
}

func TestPriority(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Task"})

	out, err := Priority(s, []string{"1", "A"})
	if err != nil {
		t.Fatalf("Priority: %v", err)
	}
	if !strings.Contains(out, "A") {
		t.Errorf("output should mention priority A, got %q", out)
	}

	tasks, _ := s.Load()
	if tasks[0].Priority != "A" {
		t.Errorf("task priority: got %q, want A", tasks[0].Priority)
	}
}

func TestPriorityRemove(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Task", "pri:A"})

	_, err := Priority(s, []string{"1", ""})
	if err != nil {
		t.Fatalf("Priority: %v", err)
	}
	tasks, _ := s.Load()
	if tasks[0].Priority != "" {
		t.Errorf("priority should be removed, got %q", tasks[0].Priority)
	}
}

func TestPriorityInvalid(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Task"})

	_, err := Priority(s, []string{"1", "1"})
	if err == nil {
		t.Error("Priority with non-letter value should return error")
	}
}

func TestPriorityMissingArgs(t *testing.T) {
	s := newTestStore(t)
	_, err := Priority(s, []string{"1"})
	if err == nil {
		t.Error("Priority with missing priority should return error")
	}
}

func TestDepri(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Task", "pri:B"})

	_, err := Depri(s, []string{"1"})
	if err != nil {
		t.Fatalf("Depri: %v", err)
	}
	tasks, _ := s.Load()
	if tasks[0].Priority != "" {
		t.Errorf("priority should be removed, got %q", tasks[0].Priority)
	}
}

func TestDepriNoPriority(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Task"})

	out, err := Depri(s, []string{"1"})
	if err != nil {
		t.Fatalf("Depri: %v", err)
	}
	if !strings.Contains(out, "sem prioridade") {
		t.Errorf("expected 'sem prioridade' message, got %q", out)
	}
}

func TestDeleteOne(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"A"})
	_, _ = Add(s, []string{"B"})
	_, _ = Add(s, []string{"C"})

	out, err := Delete(s, []string{"2"})
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !strings.Contains(out, "1 tarefa") {
		t.Errorf("expected 1 removed, got %q", out)
	}

	tasks, _ := s.Load()
	if len(tasks) != 2 {
		t.Errorf("expected 2 remaining, got %d", len(tasks))
	}
	if tasks[0].Description != "A" || tasks[1].Description != "C" {
		t.Errorf("wrong tasks remain: %v", tasks)
	}
}

func TestDeleteMultiple(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"A"})
	_, _ = Add(s, []string{"B"})
	_, _ = Add(s, []string{"C"})

	_, err := Delete(s, []string{"1", "3"})
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	tasks, _ := s.Load()
	if len(tasks) != 1 || tasks[0].Description != "B" {
		t.Errorf("expected only B to remain, got %v", tasks)
	}
}

func TestDeleteInvalidNumber(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"A"})

	_, err := Delete(s, []string{"99"})
	if err == nil {
		t.Error("Delete with invalid number should return error")
	}
}

func TestDeleteNoArgs(t *testing.T) {
	s := newTestStore(t)
	_, err := Delete(s, nil)
	if err == nil {
		t.Error("Delete with no args should return error")
	}
}

func TestArchive(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Pendente"})
	_, _ = Add(s, []string{"Concluída"})
	_, _ = Do(s, []string{"2"})

	out, err := Archive(s, nil)
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if !strings.Contains(out, "1 tarefa") {
		t.Errorf("expected 1 archived, got %q", out)
	}
	if !strings.Contains(out, s.DoneFile) {
		t.Errorf("output should mention done file, got %q", out)
	}

	remaining, _ := s.Load()
	if len(remaining) != 1 || remaining[0].Description != "Pendente" {
		t.Errorf("only pending should remain, got %v", remaining)
	}
}

func TestArchiveNoCompleted(t *testing.T) {
	s := newTestStore(t)
	_, _ = Add(s, []string{"Pendente"})

	out, err := Archive(s, nil)
	if err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if !strings.Contains(out, "Nenhuma") {
		t.Errorf("expected 'Nenhuma', got %q", out)
	}
}

func TestHelp(t *testing.T) {
	help := Help()
	for _, want := range []string{"add", "list", "do", "pri", "del", "archive", "todo.txt"} {
		if !strings.Contains(help, want) {
			t.Errorf("Help should mention %q", want)
		}
	}
}

func TestEndToEnd(t *testing.T) {
	s := newTestStore(t)

	steps := []struct {
		name string
		fn   func() (string, error)
	}{
		{"add 1", func() (string, error) { return Add(s, []string{"Tarefa A +Proj1 @ctx1"}) }},
		{"add 2", func() (string, error) { return Add(s, []string{"Tarefa B +Proj2 pri:B due:2099-01-01"}) }},
		{"add 3", func() (string, error) { return Add(s, []string{"Tarefa C +Proj1"}) }},
		{"list", func() (string, error) { return List(s, nil) }},
		{"do 1", func() (string, error) { return Do(s, []string{"1"}) }},
		{"do 2", func() (string, error) { return Do(s, []string{"2"}) }},
		{"pri 3 A", func() (string, error) { return Priority(s, []string{"3", "A"}) }},
		{"list all", func() (string, error) { return List(s, []string{"all"}) }},
		{"archive", func() (string, error) { return Archive(s, nil) }},
	}

	for _, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			out, err := step.fn()
			if err != nil {
				t.Fatalf("%s: %v", step.name, err)
			}
			if out == "" {
				t.Errorf("%s: empty output", step.name)
			}
		})
	}

	tasks, _ := s.Load()
	if len(tasks) != 1 {
		t.Errorf("expected 1 task remaining after archive, got %d", len(tasks))
	}
	if len(tasks) > 0 && tasks[0].Description != "Tarefa C +Proj1" {
		t.Errorf("wrong task remaining: %v", tasks)
	}
}
