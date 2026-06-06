package cmd

import (
	"fmt"
	"sort"
	"strings"

	"todotxt/internal/filter"
	"todotxt/internal/store"
	"todotxt/internal/task"
)

func List(s *store.Store, args []string) (string, error) {
	f := filter.Parse(args)
	if !f.OnlyCompleted && !f.ShowAll {
		f.HideCompleted = true
	}
	tasks, err := loadTasks(s)
	if err != nil {
		return "", fmt.Errorf("ao carregar tarefas: %w", err)
	}
	if len(tasks) == 0 {
		return "Nenhuma tarefa encontrada.", nil
	}

	tasks = filterAndSort(tasks, f)
	if len(tasks) == 0 {
		return "Nenhuma tarefa corresponde ao filtro.", nil
	}

	header := "Tarefas"
	if len(f.Projects) > 0 || len(f.Contexts) > 0 || f.SearchText != "" {
		header += " (filtrado)"
	}

	var sb strings.Builder
	sb.WriteString(colorize(colorBold, header) + "\n")
	sb.WriteString(strings.Repeat("─", 60) + "\n")

	for i, t := range tasks {
		sb.WriteString(formatTaskLine(i+1, t) + "\n")
	}

	sb.WriteString(fmt.Sprintf("\nTotal: %d tarefa(s)", len(tasks)))
	return sb.String(), nil
}

func filterAndSort(tasks []*task.Task, f filter.Filter) []*task.Task {
	var filtered []*task.Task
	for _, t := range tasks {
		if f.Match(t) {
			filtered = append(filtered, t)
		}
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		a, b := filtered[i], filtered[j]
		if a.Completed != b.Completed {
			return !a.Completed
		}
		if a.Priority != b.Priority {
			return priorityOrder(a.Priority) < priorityOrder(b.Priority)
		}
		da, aok := a.DueDate()
		db, bok := b.DueDate()
		if aok != bok {
			return aok
		}
		if aok && bok && !da.Equal(db) {
			return da.Before(db)
		}
		return a.Description < b.Description
	})

	return filtered
}

func priorityOrder(p string) int {
	if p == "" {
		return 100
	}
	return int(p[0] - 'A')
}
