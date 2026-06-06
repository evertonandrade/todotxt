package task

import (
	"strings"
	"testing"
	"time"
)

func TestParseSimple(t *testing.T) {
	tk := Parse("Buy milk", 1)
	if tk.Description != "Buy milk" {
		t.Errorf("Description: got %q, want %q", tk.Description, "Buy milk")
	}
	if tk.Completed || tk.Priority != "" {
		t.Errorf("expected pending+no-priority, got completed=%v priority=%q", tk.Completed, tk.Priority)
	}
}

func TestParseWithPriority(t *testing.T) {
	tk := Parse("(A) Call Mom +Family @phone", 1)
	if tk.Priority != "A" {
		t.Errorf("Priority: got %q, want A", tk.Priority)
	}
	if tk.Description != "Call Mom +Family @phone" {
		t.Errorf("Description: got %q", tk.Description)
	}
	if !tk.HasProject("Family") {
		t.Error("should have project Family")
	}
	if !tk.HasContext("phone") {
		t.Error("should have context phone")
	}
}

func TestParseWithCreationDate(t *testing.T) {
	tk := Parse("2024-02-01 Buy groceries", 1)
	if tk.CreationDate != "2024-02-01" {
		t.Errorf("CreationDate: got %q, want 2024-02-01", tk.CreationDate)
	}
	if tk.Description != "Buy groceries" {
		t.Errorf("Description: got %q", tk.Description)
	}
}

func TestParsePriorityAndDate(t *testing.T) {
	tk := Parse("(A) 2024-01-15 Call Mom +Family", 1)
	if tk.Priority != "A" {
		t.Errorf("Priority: got %q", tk.Priority)
	}
	if tk.CreationDate != "2024-01-15" {
		t.Errorf("CreationDate: got %q", tk.CreationDate)
	}
}

func TestParseCompleted(t *testing.T) {
	tk := Parse("x 2024-01-15 2024-01-10 Call Mom", 1)
	if !tk.Completed {
		t.Error("should be completed")
	}
	if tk.CompletionDate != "2024-01-15" {
		t.Errorf("CompletionDate: got %q", tk.CompletionDate)
	}
	if tk.CreationDate != "2024-01-10" {
		t.Errorf("CreationDate: got %q", tk.CreationDate)
	}
	if tk.Description != "Call Mom" {
		t.Errorf("Description: got %q", tk.Description)
	}
}

func TestParseCompletedSimple(t *testing.T) {
	tk := Parse("x Review PR", 1)
	if !tk.Completed {
		t.Error("should be completed")
	}
	if tk.CompletionDate != "" {
		t.Errorf("CompletionDate should be empty, got %q", tk.CompletionDate)
	}
	if tk.Description != "Review PR" {
		t.Errorf("Description: got %q", tk.Description)
	}
}

func TestParseCompletedWithPriority(t *testing.T) {
	tk := Parse("x 2024-01-15 (B) 2024-01-10 Call Mom", 1)
	if !tk.Completed {
		t.Error("should be completed")
	}
	if tk.Priority != "B" {
		t.Errorf("Priority: got %q", tk.Priority)
	}
	if tk.CompletionDate != "2024-01-15" {
		t.Errorf("CompletionDate: got %q", tk.CompletionDate)
	}
	if tk.CreationDate != "2024-01-10" {
		t.Errorf("CreationDate: got %q", tk.CreationDate)
	}
}

func TestParseEmptyLine(t *testing.T) {
	tk := Parse("", 1)
	if tk.Description != "" {
		t.Errorf("expected empty description, got %q", tk.Description)
	}
}

func TestParseMultipleProjectsAndContexts(t *testing.T) {
	tk := Parse("Reunião +Trabalho +Urgente @escritório @reunião", 1)
	if len(tk.Projects) != 2 {
		t.Errorf("expected 2 projects, got %v", tk.Projects)
	}
	if len(tk.Contexts) != 2 {
		t.Errorf("expected 2 contexts, got %v", tk.Contexts)
	}
}

func TestParseCustomFields(t *testing.T) {
	tk := Parse("Task due:2024-12-31 t:2024-12-01 key:value", 1)
	if v, ok := tk.Custom("due"); !ok || v != "2024-12-31" {
		t.Errorf("due: got %q, ok=%v", v, ok)
	}
	if v, ok := tk.Custom("t"); !ok || v != "2024-12-01" {
		t.Errorf("t: got %q, ok=%v", v, ok)
	}
	if v, ok := tk.Custom("key"); !ok || v != "value" {
		t.Errorf("key: got %q, ok=%v", v, ok)
	}
}

func TestParseURLNotCustomField(t *testing.T) {
	tk := Parse("Ver https://example.com/path:8080 docs", 1)
	if _, ok := tk.Custom("path"); ok {
		t.Error("URL with :port should not be parsed as custom field")
	}
}

func TestParseCaseInsensitivePriority(t *testing.T) {
	tk := Parse("(a) Low priority", 1)
	if tk.Priority != "a" {
		t.Errorf("Priority: got %q", tk.Priority)
	}
}

func TestFormatRoundtrip(t *testing.T) {
	lines := []string{
		"Buy milk",
		"(A) Call Mom +Family @phone",
		"2024-02-01 Buy groceries",
		"(A) 2024-01-15 Call Mom +Family @phone",
		"x 2024-01-15 2024-01-10 Call Mom +Family @phone",
		"x Review PR",
		"Task due:2024-12-31 +Project @context",
	}
	for _, line := range lines {
		t.Run(line, func(t *testing.T) {
			first := Parse(line, 1)
			out := first.Format()
			second := Parse(out, 1)
			if first.Priority != second.Priority {
				t.Errorf("Priority lost: %q -> %q", first.Priority, second.Priority)
			}
			if first.Completed != second.Completed {
				t.Errorf("Completed lost: %v -> %v", first.Completed, second.Completed)
			}
			if first.CreationDate != second.CreationDate {
				t.Errorf("CreationDate lost: %q -> %q", first.CreationDate, second.CreationDate)
			}
			if first.CompletionDate != second.CompletionDate {
				t.Errorf("CompletionDate lost: %q -> %q", first.CompletionDate, second.CompletionDate)
			}
			if first.Description != second.Description {
				t.Errorf("Description lost: %q -> %q", first.Description, second.Description)
			}
		})
	}
}

func TestMarkComplete(t *testing.T) {
	tk := Parse("Buy milk", 1)
	if tk.Completed {
		t.Fatal("should not be completed initially")
	}
	tk.MarkComplete()
	if !tk.Completed {
		t.Error("should be completed after MarkComplete")
	}
	date, err := time.Parse("2006-01-02", tk.CompletionDate)
	if err != nil {
		t.Fatalf("CompletionDate invalid: %q", tk.CompletionDate)
	}
	today := time.Now()
	if date.Year() != today.Year() || date.YearDay() != today.YearDay() {
		t.Errorf("CompletionDate: got %s, expected today (%s)", tk.CompletionDate, today.Format("2006-01-02"))
	}
}

func TestMarkIncomplete(t *testing.T) {
	tk := Parse("x 2024-01-01 Done task", 1)
	tk.MarkIncomplete()
	if tk.Completed {
		t.Error("should not be completed after MarkIncomplete")
	}
	if tk.CompletionDate != "" {
		t.Errorf("CompletionDate should be empty, got %q", tk.CompletionDate)
	}
}

func TestSetPriority(t *testing.T) {
	tk := Parse("Buy milk", 1)
	tk.SetPriority("A")
	if tk.Priority != "A" {
		t.Errorf("Priority: got %q, want A", tk.Priority)
	}
	tk.SetPriority("b")
	if tk.Priority != "B" {
		t.Errorf("Priority should be uppercase, got %q", tk.Priority)
	}
	tk.SetPriority("")
	if tk.Priority != "" {
		t.Errorf("Priority should be empty, got %q", tk.Priority)
	}
}

func TestHasProjectCaseInsensitive(t *testing.T) {
	tk := Parse("Task +Trabalho", 1)
	if !tk.HasProject("trabalho") {
		t.Error("HasProject should be case-insensitive")
	}
	if !tk.HasProject("TRABALHO") {
		t.Error("HasProject should be case-insensitive")
	}
}

func TestDueDate(t *testing.T) {
	tk := Parse("Task due:2024-12-31", 1)
	due, ok := tk.DueDate()
	if !ok {
		t.Fatal("DueDate should be present")
	}
	if due.Year() != 2024 || due.Month() != 12 || due.Day() != 31 {
		t.Errorf("DueDate: got %s, want 2024-12-31", due.Format("2006-01-02"))
	}
}

func TestDueDateAbsent(t *testing.T) {
	tk := Parse("Task without due", 1)
	if _, ok := tk.DueDate(); ok {
		t.Error("DueDate should not be present")
	}
}

func TestIsOverdue(t *testing.T) {
	past := Parse("Task due:2020-01-01", 1)
	if !past.IsOverdue() {
		t.Error("past due date should be overdue")
	}
	if past.IsDueToday() {
		t.Error("past due date is not today")
	}

	future := Parse("Task due:2099-12-31", 1)
	if future.IsOverdue() {
		t.Error("future due date should not be overdue")
	}

	noDue := Parse("Task no due", 1)
	if noDue.IsOverdue() {
		t.Error("task without due date should not be overdue")
	}

	completed := Parse("x 2020-01-01 Task done due:2020-01-01", 1)
	if completed.IsOverdue() {
		t.Error("completed task should not be overdue")
	}
}

func TestIsDueSoon(t *testing.T) {
	future := time.Now().AddDate(0, 0, 3).Format("2006-01-02")
	tk := Parse("Task due:"+future, 1)
	if !tk.IsDueSoon(7) {
		t.Error("task in 3 days should be due soon (within 7 days)")
	}
	if tk.IsDueSoon(2) {
		t.Error("task in 3 days should not be due soon (within 2 days)")
	}
}

func TestNewConstructor(t *testing.T) {
	tk := New("Task +proj @ctx due:2024-12-31")
	if tk.Description != "Task +proj @ctx due:2024-12-31" {
		t.Errorf("Description: got %q", tk.Description)
	}
	if len(tk.Projects) != 1 || tk.Projects[0] != "proj" {
		t.Errorf("Projects: %v", tk.Projects)
	}
	if len(tk.Contexts) != 1 || tk.Contexts[0] != "ctx" {
		t.Errorf("Contexts: %v", tk.Contexts)
	}
	if v, ok := tk.Custom("due"); !ok || v != "2024-12-31" {
		t.Errorf("due: got %q, ok=%v", v, ok)
	}
}

func TestLineNumberPreserved(t *testing.T) {
	tk := Parse("Task", 42)
	if tk.LineNumber != 42 {
		t.Errorf("LineNumber: got %d, want 42", tk.LineNumber)
	}
}

func TestUniqueProjects(t *testing.T) {
	tk := Parse("Task +proj +proj +outro", 1)
	count := 0
	for _, p := range tk.Projects {
		if p == "proj" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected unique projects, got %v", tk.Projects)
	}
}

func TestFormatEmptyDescription(t *testing.T) {
	tk := &Task{Description: "Test"}
	out := tk.Format()
	if !strings.Contains(out, "Test") {
		t.Errorf("Format should include description, got %q", out)
	}
}
