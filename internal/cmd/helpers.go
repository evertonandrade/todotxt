package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"todotxt/internal/store"
	"todotxt/internal/task"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

var noColor = false

func init() {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		noColor = true
	}
}

func colorize(c, s string) string {
	if noColor {
		return s
	}
	return c + s + colorReset
}

func priorityColor(p string) string {
	switch strings.ToUpper(p) {
	case "A":
		return colorRed
	case "B":
		return colorYellow
	case "C":
		return colorBlue
	default:
		return colorGray
	}
}

func formatNumber(n int) string {
	return colorize(colorCyan, fmt.Sprintf("%3d", n))
}

func formatTaskLine(n int, t *task.Task) string {
	marker := " "
	if t.Completed {
		marker = "x"
	}

	var sb strings.Builder
	sb.WriteString(formatNumber(n))
	sb.WriteString(" ")
	sb.WriteString(colorize(colorGray, "["))
	sb.WriteString(colorize(colorGray, marker))
	sb.WriteString(colorize(colorGray, "] "))

	if t.Priority != "" && !t.Completed {
		sb.WriteString(colorize(priorityColor(t.Priority), "("+t.Priority+")"))
		sb.WriteString(" ")
	}

	desc := t.Description

	dueStr := ""
	if !t.Completed {
		if due, ok := t.DueDate(); ok {
			switch {
			case t.IsOverdue():
				dueStr = " " + colorize(colorRed+colorBold, fmt.Sprintf("[venceu: %s]", due.Format("2006-01-02")))
			case t.IsDueToday():
				dueStr = " " + colorize(colorRed, "[vence hoje]")
			default:
				days := int(due.Sub(time.Now()).Hours() / 24)
				if days <= 7 {
					dueStr = " " + colorize(colorYellow, fmt.Sprintf("[vence em %dd: %s]", days, due.Format("2006-01-02")))
				} else {
					dueStr = " " + colorize(colorGray, fmt.Sprintf("[vence: %s]", due.Format("2006-01-02")))
				}
			}
		}
	}

	if t.Completed {
		desc = colorize(colorGray+colorBold, "(x) ") + colorize(colorGray, desc)
		if t.CompletionDate != "" {
			desc += " " + colorize(colorGray, "["+t.CompletionDate+"]")
		}
	} else {
		desc = colorizeProjectContext(desc)
	}

	sb.WriteString(desc)
	sb.WriteString(dueStr)
	return sb.String()
}

func colorizeProjectContext(s string) string {
	var sb strings.Builder
	words := strings.Fields(s)
	for i, w := range words {
		if i > 0 {
			sb.WriteString(" ")
		}
		switch {
		case strings.HasPrefix(w, "+") && len(w) > 1:
			sb.WriteString(colorize(colorGreen, w))
		case strings.HasPrefix(w, "@") && len(w) > 1:
			sb.WriteString(colorize(colorCyan, w))
		case strings.Contains(w, ":") && !strings.HasPrefix(w, "http"):
			idx := strings.Index(w, ":")
			if idx > 0 && idx < len(w)-1 {
				key := w[:idx]
				if isAlphaNumericKey(key) && !isDateValue(w[idx+1:]) {
					sb.WriteString(colorize(colorYellow, w))
					continue
				}
			}
			sb.WriteString(w)
		default:
			sb.WriteString(w)
		}
	}
	return sb.String()
}

func isAlphaNumericKey(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}

func isDateValue(s string) bool {
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}

func loadTasks(s *store.Store) []*task.Task {
	tasks, err := s.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao carregar tarefas: %v\n", err)
		os.Exit(1)
	}
	return tasks
}

func saveTasks(s *store.Store, tasks []*task.Task) {
	if err := s.Save(tasks); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao guardar tarefas: %v\n", err)
		os.Exit(1)
	}
}

func printErr(w io.Writer, msg string) {
	fmt.Fprintln(w, colorize(colorRed, "✗ ")+msg)
}

func printOK(msg string) {
	fmt.Println(colorize(colorGreen, "✓ ") + msg)
}
