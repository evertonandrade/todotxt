package task

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Completed      bool
	CompletionDate string
	Priority       string
	CreationDate   string
	Description    string
	Projects       []string
	Contexts       []string
	CustomFields   map[string]string
	Raw            string
	LineNumber     int
}

func New(description string) *Task {
	t := &Task{
		Description:  description,
		CustomFields: make(map[string]string),
	}
	t.extractMetadata()
	return t
}

func Parse(line string, lineNumber int) *Task {
	t := &Task{
		Raw:          line,
		CustomFields: make(map[string]string),
		LineNumber:   lineNumber,
	}

	rest := strings.TrimSpace(line)
	if rest == "" {
		return t
	}

	fields := strings.Fields(rest)
	if len(fields) == 0 {
		return t
	}

	i := 0

	if fields[i] == "x" {
		t.Completed = true
		i++
		if i < len(fields) && isDate(fields[i]) {
			t.CompletionDate = fields[i]
			i++
		}
	}

	if i < len(fields) && isPriority(fields[i]) {
		t.Priority = fields[i][1:2]
		i++
	}

	if !t.Completed && i < len(fields) && isDate(fields[i]) {
		t.CreationDate = fields[i]
		i++
	} else if t.Completed && i < len(fields) && isDate(fields[i]) {
		t.CreationDate = fields[i]
		i++
	}

	t.Description = strings.Join(fields[i:], " ")
	t.extractMetadata()
	return t
}

func (t *Task) extractMetadata() {
	t.Projects = nil
	t.Contexts = nil
	t.CustomFields = make(map[string]string)

	words := strings.Fields(t.Description)
	for _, word := range words {
		switch {
		case strings.HasPrefix(word, "+") && len(word) > 1 && !strings.ContainsAny(word[1:], " \t"):
			t.Projects = appendUnique(t.Projects, word[1:])
		case strings.HasPrefix(word, "@") && len(word) > 1 && !strings.ContainsAny(word[1:], " \t"):
			t.Contexts = appendUnique(t.Contexts, word[1:])
		default:
			if idx := strings.Index(word, ":"); idx > 0 && idx < len(word)-1 {
				key := word[:idx]
				value := word[idx+1:]
				if isAlpha(key) && !strings.Contains(key, "/") && !strings.HasPrefix(word, "http") {
					t.CustomFields[key] = value
				}
			}
		}
	}
}

func (t *Task) HasProject(name string) bool {
	for _, p := range t.Projects {
		if strings.EqualFold(p, name) {
			return true
		}
	}
	return false
}

func (t *Task) HasContext(name string) bool {
	for _, c := range t.Contexts {
		if strings.EqualFold(c, name) {
			return true
		}
	}
	return false
}

func (t *Task) Custom(key string) (string, bool) {
	v, ok := t.CustomFields[strings.ToLower(key)]
	return v, ok
}

func (t *Task) DueDate() (time.Time, bool) {
	if v, ok := t.Custom("due"); ok {
		return parseDate(v)
	}
	return time.Time{}, false
}

func (t *Task) IsOverdue() bool {
	if t.Completed {
		return false
	}
	due, ok := t.DueDate()
	if !ok {
		return false
	}
	return time.Now().After(due)
}

func (t *Task) IsDueToday() bool {
	due, ok := t.DueDate()
	if !ok {
		return false
	}
	return sameDay(due, time.Now())
}

func (t *Task) IsDueSoon(days int) bool {
	due, ok := t.DueDate()
	if !ok {
		return false
	}
	now := time.Now()
	if due.Before(now) {
		return false
	}
	return due.Sub(now).Hours() <= float64(days*24)
}

func (t *Task) Format() string {
	var parts []string

	if t.Completed {
		parts = append(parts, "x")
		if t.CompletionDate != "" {
			parts = append(parts, t.CompletionDate)
		}
		if t.CreationDate != "" {
			parts = append(parts, t.CreationDate)
		}
		parts = append(parts, t.Description)
	} else {
		if t.Priority != "" {
			parts = append(parts, "("+strings.ToUpper(t.Priority)+")")
		}
		if t.CreationDate != "" {
			parts = append(parts, t.CreationDate)
		}
		parts = append(parts, t.Description)
	}
	return strings.Join(parts, " ")
}

func (t *Task) MarkComplete() {
	t.Completed = true
	t.CompletionDate = time.Now().Format("2006-01-02")
}

func (t *Task) MarkIncomplete() {
	t.Completed = false
	t.CompletionDate = ""
}

func (t *Task) SetPriority(p string) {
	if p == "" {
		t.Priority = ""
		return
	}
	p = strings.ToUpper(strings.TrimSpace(p))
	if len(p) > 1 {
		p = string(p[0])
	}
	t.Priority = p
}

func isDate(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func parseDate(s string) (time.Time, bool) {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, false
	}
	return d, true
}

func isPriority(s string) bool {
	if len(s) != 3 {
		return false
	}
	if s[0] != '(' || s[2] != ')' {
		return false
	}
	return isAlpha(string(s[1]))
}

func isAlpha(s string) bool {
	re := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	return re.MatchString(s)
}

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func appendUnique(slice []string, val string) []string {
	for _, v := range slice {
		if v == val {
			return slice
		}
	}
	return append(slice, val)
}

func ParsePriority(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}
	s = strings.ToUpper(s)
	if len(s) > 1 {
		s = string(s[0])
	}
	if !isAlpha(s) {
		return "", fmt.Errorf("prioridade inválida: %q", s)
	}
	return s, nil
}

func ParseLineNumber(s string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, fmt.Errorf("número inválido: %q", s)
	}
	if n <= 0 {
		return 0, fmt.Errorf("número deve ser positivo: %d", n)
	}
	return n, nil
}
