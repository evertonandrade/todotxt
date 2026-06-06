package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"todotxt/internal/store"
	"todotxt/internal/task"
)

func Add(s *store.Store, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Erro: forneça uma descrição para a tarefa.")
		fmt.Fprintln(os.Stderr, "Uso: todotxt add <descrição> [+projeto] [@contexto] [pri:A] [due:YYYY-MM-DD]")
		os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "Erro: data inválida para due: %q (use YYYY-MM-DD)\n", rest)
			os.Exit(1)
		}
		due = rest
		raw = stripToken(raw, "due:")
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		fmt.Fprintln(os.Stderr, "Erro: descrição vazia.")
		os.Exit(1)
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

	tasks := loadTasks(s)
	tasks = append(tasks, task.Parse(line, len(tasks)+1))
	saveTasks(s, tasks)

	printOK("Tarefa adicionada:")
	fmt.Println("  " + line)
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
