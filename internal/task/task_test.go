package task

import (
	"testing"
)

func TestParseSimple(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		priority   string
		completed  bool
		creation   string
		completion string
		desc       string
	}{
		{
			name:  "simple task",
			input: "Buy milk",
			desc:  "Buy milk",
		},
		{
			name:     "with priority",
			input:    "(A) Call Mom +Family @phone",
			priority: "A",
			desc:     "Call Mom +Family @phone",
		},
		{
			name:     "with priority and creation date",
			input:    "(A) 2024-01-15 Call Mom +Family @phone",
			priority: "A",
			creation: "2024-01-15",
			desc:     "Call Mom +Family @phone",
		},
		{
			name:       "completed task",
			input:      "x 2024-01-15 2024-01-10 Call Mom +Family @phone",
			completed:  true,
			completion: "2024-01-15",
			creation:   "2024-01-10",
			desc:       "Call Mom +Family @phone",
		},
		{
			name:       "completed with priority",
			input:      "x 2024-01-15 (B) 2024-01-10 Call Mom +Family @phone",
			completed:  true,
			priority:   "B",
			completion: "2024-01-15",
			creation:   "2024-01-10",
			desc:       "Call Mom +Family @phone",
		},
		{
			name:     "creation date only",
			input:    "2024-02-01 Buy groceries",
			creation: "2024-02-01",
			desc:     "Buy groceries",
		},
		{
			name:      "completed simple",
			input:     "x Review PR",
			completed: true,
			desc:      "Review PR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Parse(tt.input, 1)
			if task.Completed != tt.completed {
				t.Errorf("Completed: got %v, want %v", task.Completed, tt.completed)
			}
			if task.Priority != tt.priority {
				t.Errorf("Priority: got %q, want %q", task.Priority, tt.priority)
			}
			if task.CreationDate != tt.creation {
				t.Errorf("CreationDate: got %q, want %q", task.CreationDate, tt.creation)
			}
			if task.CompletionDate != tt.completion {
				t.Errorf("CompletionDate: got %q, want %q", task.CompletionDate, tt.completion)
			}
			if task.Description != tt.desc {
				t.Errorf("Description: got %q, want %q", task.Description, tt.desc)
			}
		})
	}
}

func TestExtractProjectsAndContexts(t *testing.T) {
	task := Parse("Buy milk +Shopping @store due:2024-02-05", 1)
	if len(task.Projects) != 1 || task.Projects[0] != "Shopping" {
		t.Errorf("Projects: got %v, want [Shopping]", task.Projects)
	}
	if len(task.Contexts) != 1 || task.Contexts[0] != "store" {
		t.Errorf("Contexts: got %v, want [store]", task.Contexts)
	}
	if v, ok := task.Custom("due"); !ok || v != "2024-02-05" {
		t.Errorf("due custom field: got %q, want 2024-02-05", v)
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"simple", "Buy milk"},
		{"with priority", "(A) Call Mom +Family"},
		{"with creation date", "2024-01-15 Call Mom +Family"},
		{"completed", "x 2024-01-15 2024-01-10 Call Mom +Family"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Parse(tt.input, 1)
			formatted := task.Format()
			parsed := Parse(formatted, 1)
			if parsed.Priority != task.Priority {
				t.Errorf("Priority mismatch after roundtrip: %q vs %q", parsed.Priority, task.Priority)
			}
			if parsed.Completed != task.Completed {
				t.Errorf("Completed mismatch after roundtrip: %v vs %v", parsed.Completed, task.Completed)
			}
			if parsed.Description != task.Description {
				t.Errorf("Description mismatch: %q vs %q", parsed.Description, task.Description)
			}
		})
	}
}

func TestMarkComplete(t *testing.T) {
	task := Parse("Buy milk", 1)
	if task.Completed {
		t.Error("should not be completed initially")
	}
	task.MarkComplete()
	if !task.Completed {
		t.Error("should be completed after MarkComplete")
	}
	if task.CompletionDate == "" {
		t.Error("CompletionDate should be set")
	}
}

func TestSetPriority(t *testing.T) {
	task := Parse("Buy milk", 1)
	task.SetPriority("A")
	if task.Priority != "A" {
		t.Errorf("Priority: got %q, want A", task.Priority)
	}
	task.SetPriority("")
	if task.Priority != "" {
		t.Errorf("Priority: got %q, want empty", task.Priority)
	}
}
