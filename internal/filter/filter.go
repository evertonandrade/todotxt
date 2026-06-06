package filter

import (
	"strings"

	"todotxt/internal/task"
)

type Filter struct {
	Projects      []string
	Contexts      []string
	Priorities    []string
	HideCompleted bool
	OnlyCompleted bool
	ShowAll       bool
	SearchText    string
	HasDueDate    bool
	Overdue       bool
	DueToday      bool
	DueSoonDays   int
}

func Parse(args []string) Filter {
	f := Filter{}
	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "+") && len(a) > 1:
			f.Projects = append(f.Projects, a[1:])
		case strings.HasPrefix(a, "@") && len(a) > 1:
			f.Contexts = append(f.Contexts, a[1:])
		case strings.HasPrefix(a, "pri:"):
			f.Priorities = append(f.Priorities, strings.ToUpper(a[4:]))
		case a == "-x" || a == "--hide-completed":
			f.HideCompleted = true
		case a == "done" || a == "--completed":
			f.OnlyCompleted = true
		case a == "all":
			f.ShowAll = true
		case a == "due" || a == "--has-due":
			f.HasDueDate = true
		case a == "overdue":
			f.Overdue = true
		case a == "today":
			f.DueToday = true
		case strings.HasPrefix(a, "due:"):
			f.HasDueDate = true
		default:
			if f.SearchText != "" {
				f.SearchText += " "
			}
			f.SearchText += a
		}
	}
	return f
}

func (f Filter) Match(t *task.Task) bool {
	if f.HideCompleted && t.Completed {
		return false
	}
	if f.OnlyCompleted && !t.Completed {
		return false
	}

	for _, p := range f.Projects {
		if !t.HasProject(p) {
			return false
		}
	}
	for _, c := range f.Contexts {
		if !t.HasContext(c) {
			return false
		}
	}
	for _, p := range f.Priorities {
		if !strings.EqualFold(t.Priority, p) {
			return false
		}
	}

	if f.HasDueDate {
		if _, ok := t.DueDate(); !ok {
			return false
		}
	}
	if f.Overdue && !t.IsOverdue() {
		return false
	}
	if f.DueToday && !t.IsDueToday() {
		return false
	}
	if f.DueSoonDays > 0 && !t.IsDueSoon(f.DueSoonDays) {
		return false
	}

	if f.SearchText != "" {
		if !strings.Contains(strings.ToLower(t.Description), strings.ToLower(f.SearchText)) {
			return false
		}
	}

	return true
}
