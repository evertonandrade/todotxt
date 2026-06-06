package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"todotxt/internal/store"
	"todotxt/internal/task"
)

func Add(s *store.Store, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("forneça uma descrição para a tarefa. Uso: todotxt add <descrição> [+projeto] [@contexto] [pri:A] [due:YYYY-MM-DD]")
	}

	raw := strings.Join(args, " ")

	priority := ""
	if rest, ok := extractToken(raw, "pri:"); ok {
		priority = strings.ToUpper(rest)
		raw = stripToken(raw, "pri:")
	}

	due := ""
	if rest, ok := extractToken(raw, "due:"); ok {
		if _, err := time.Parse("2006-01-02", rest); err != nil {
			return "", fmt.Errorf("data inválida para due: %q (use YYYY-MM-DD)", rest)
		}
		due = rest
		raw = stripToken(raw, "due:")
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("descrição vazia")
	}

	parts := []string{}
	if priority != "" {
		parts = append(parts, "("+priority+")")
	}
	parts = append(parts, time.Now().Format("2006-01-02"))
	parts = append(parts, raw)
	if due != "" {
		parts = append(parts, "due:"+due)
	}
	line := strings.Join(parts, " ")

	tasks, err := loadTasks(s)
	if err != nil {
		return "", fmt.Errorf("ao carregar tarefas: %w", err)
	}
	tasks = append(tasks, task.Parse(line, len(tasks)+1))
	if err := saveTasks(s, tasks); err != nil {
		return "", fmt.Errorf("ao guardar tarefas: %w", err)
	}

	return fmt.Sprintf("Tarefa adicionada:\n  %s", line), nil
}

func extractToken(s, key string) (string, bool) {
	idx := strings.Index(s, " "+key)
	if idx < 0 {
		if strings.HasPrefix(s, key) {
			return "", true
		}
		return "", false
	}
	rest := s[idx+len(" ")+len(key):]
	if sp := strings.IndexAny(rest, " \t"); sp >= 0 {
		rest = rest[:sp]
	}
	return rest, true
}

func stripToken(s, key string) string {
	idx := strings.Index(s, " "+key)
	if idx < 0 {
		return s
	}
	rest := s[idx+len(" ")+len(key):]
	if sp := strings.IndexAny(rest, " \t"); sp >= 0 {
		rest = rest[sp:]
	} else {
		rest = ""
	}
	return strings.TrimSpace(s[:idx] + rest)
}
