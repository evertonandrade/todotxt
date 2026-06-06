package filter

import (
	"testing"

	"todotxt/internal/task"
)

func mkTask(desc string) *task.Task {
	return task.New(desc)
}

func TestParseEmpty(t *testing.T) {
	f := Parse(nil)
	if len(f.Projects) != 0 || len(f.Contexts) != 0 {
		t.Errorf("expected empty filter, got %+v", f)
	}
}

func TestParseProjects(t *testing.T) {
	f := Parse([]string{"+Trabalho", "+Urgente"})
	if len(f.Projects) != 2 || f.Projects[0] != "Trabalho" || f.Projects[1] != "Urgente" {
		t.Errorf("Projects: got %v, want [Trabalho Urgente]", f.Projects)
	}
}

func TestParseContexts(t *testing.T) {
	f := Parse([]string{"@casa", "@escritório"})
	if len(f.Contexts) != 2 || f.Contexts[0] != "casa" || f.Contexts[1] != "escritório" {
		t.Errorf("Contexts: got %v, want [casa escritório]", f.Contexts)
	}
}

func TestParsePriorities(t *testing.T) {
	f := Parse([]string{"pri:A", "pri:B"})
	if len(f.Priorities) != 2 || f.Priorities[0] != "A" || f.Priorities[1] != "B" {
		t.Errorf("Priorities: got %v, want [A B]", f.Priorities)
	}
}

func TestParseShowFlags(t *testing.T) {
	tests := []struct {
		args     []string
		hide     bool
		onlyDone bool
		showAll  bool
		hasDue   bool
		overdue  bool
		dueToday bool
	}{
		{nil, false, false, false, false, false, false},
		{[]string{"-x"}, true, false, false, false, false, false},
		{[]string{"done"}, false, true, false, false, false, false},
		{[]string{"all"}, false, false, true, false, false, false},
		{[]string{"due"}, false, false, false, true, false, false},
		{[]string{"overdue"}, false, false, false, false, true, false},
		{[]string{"today"}, false, false, false, false, false, true},
		{[]string{"--hide-completed"}, true, false, false, false, false, false},
		{[]string{"--completed"}, false, true, false, false, false, false},
		{[]string{"--has-due"}, false, false, false, true, false, false},
	}

	for _, tt := range tests {
		f := Parse(tt.args)
		if f.HideCompleted != tt.hide {
			t.Errorf("args=%v: HideCompleted got %v, want %v", tt.args, f.HideCompleted, tt.hide)
		}
		if f.OnlyCompleted != tt.onlyDone {
			t.Errorf("args=%v: OnlyCompleted got %v, want %v", tt.args, f.OnlyCompleted, tt.onlyDone)
		}
		if f.ShowAll != tt.showAll {
			t.Errorf("args=%v: ShowAll got %v, want %v", tt.args, f.ShowAll, tt.showAll)
		}
		if f.HasDueDate != tt.hasDue {
			t.Errorf("args=%v: HasDueDate got %v, want %v", tt.args, f.HasDueDate, tt.hasDue)
		}
		if f.Overdue != tt.overdue {
			t.Errorf("args=%v: Overdue got %v, want %v", tt.args, f.Overdue, tt.overdue)
		}
		if f.DueToday != tt.dueToday {
			t.Errorf("args=%v: DueToday got %v, want %v", tt.args, f.DueToday, tt.dueToday)
		}
	}
}

func TestParseSearchText(t *testing.T) {
	f := Parse([]string{"comprar", "leite"})
	if f.SearchText != "comprar leite" {
		t.Errorf("SearchText: got %q, want %q", f.SearchText, "comprar leite")
	}
}

func TestParseMixed(t *testing.T) {
	f := Parse([]string{"+Trabalho", "@escritório", "pri:A", "overdue", "reunião"})
	if len(f.Projects) != 1 || f.Projects[0] != "Trabalho" {
		t.Errorf("Projects: %v", f.Projects)
	}
	if len(f.Contexts) != 1 || f.Contexts[0] != "escritório" {
		t.Errorf("Contexts: %v", f.Contexts)
	}
	if len(f.Priorities) != 1 || f.Priorities[0] != "A" {
		t.Errorf("Priorities: %v", f.Priorities)
	}
	if !f.Overdue {
		t.Errorf("Overdue should be true")
	}
	if f.SearchText != "reunião" {
		t.Errorf("SearchText: got %q", f.SearchText)
	}
}

func TestMatchProject(t *testing.T) {
	f := Parse([]string{"+Trabalho"})
	if !f.Match(mkTask("Reunião +Trabalho @escritório")) {
		t.Error("should match +Trabalho")
	}
	if f.Match(mkTask("Comprar leite +Pessoal")) {
		t.Error("should not match +Pessoal when filtering +Trabalho")
	}
}

func TestMatchContext(t *testing.T) {
	f := Parse([]string{"@casa"})
	if !f.Match(mkTask("Lavar louça @casa")) {
		t.Error("should match @casa")
	}
	if f.Match(mkTask("Estudar @escritório")) {
		t.Error("should not match @escritório when filtering @casa")
	}
}

func TestMatchMultipleProjectsAND(t *testing.T) {
	f := Parse([]string{"+Trabalho", "+Urgente"})
	if !f.Match(mkTask("Reunião +Trabalho +Urgente")) {
		t.Error("should match task with both projects")
	}
	if f.Match(mkTask("Reunião +Trabalho")) {
		t.Error("should not match task with only one of the projects (AND logic)")
	}
}

func TestMatchPriority(t *testing.T) {
	f := Parse([]string{"pri:A"})
	a := task.New("Tarefa A")
	a.SetPriority("A")
	b := task.New("Tarefa B")
	b.SetPriority("B")

	if !f.Match(a) {
		t.Error("A should match pri:A")
	}
	if f.Match(b) {
		t.Error("B should not match pri:A")
	}
	if f.Match(mkTask("Sem prioridade")) {
		t.Error("task without priority should not match pri:A")
	}
}

func TestMatchCompletedFlags(t *testing.T) {
	hide := Parse([]string{"-x"})
	only := Parse([]string{"done"})
	all := Parse([]string{"all"})

	completed := task.New("Feito")
	completed.MarkComplete()
	pending := task.New("Pendente")

	if hide.Match(completed) {
		t.Error("-x should hide completed")
	}
	if !hide.Match(pending) {
		t.Error("-x should show pending")
	}
	if !only.Match(completed) {
		t.Error("done should show completed")
	}
	if only.Match(pending) {
		t.Error("done should hide pending")
	}
	if !all.Match(completed) || !all.Match(pending) {
		t.Error("all should match both")
	}
}

func TestMatchSearchText(t *testing.T) {
	f := Parse([]string{"comprar"})
	if !f.Match(mkTask("Comprar leite")) {
		t.Error("search should be case-insensitive substring match")
	}
	if f.Match(mkTask("Vender fruta")) {
		t.Error("should not match unrelated text")
	}
}

func TestMatchDueFlags(t *testing.T) {
	due := Parse([]string{"due"})
	if !due.Match(mkTask("Tarefa due:2026-12-31")) {
		t.Error("due should match task with due date")
	}
	if due.Match(mkTask("Tarefa sem data")) {
		t.Error("due should not match task without due date")
	}
}
