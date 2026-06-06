package cmd

import (
	"fmt"
	"sort"
	"strings"

	"todotxt/internal/filter"
	"todotxt/internal/store"
	"todotxt/internal/task"
)

func List(s *store.Store, args []string) {
	f := filter.Parse(args)
	if !f.OnlyCompleted && !f.ShowAll {
		f.HideCompleted = true
	}
	tasks := loadTasks(s)
	if len(tasks) == 0 {
		fmt.Println("Nenhuma tarefa encontrada.")
		return
	}

	tasks = filterAndSort(tasks, f)
	if len(tasks) == 0 {
		fmt.Println("Nenhuma tarefa corresponde ao filtro.")
		return
	}

	header := "Tarefas"
	if len(f.Projects) > 0 || len(f.Contexts) > 0 || f.SearchText != "" {
		header += " (filtrado)"
	}
	fmt.Println(colorize(colorBold, header))
	fmt.Println(strings.Repeat("─", 60))

	for i, t := range tasks {
		fmt.Println(formatTaskLine(i+1, t))
	}

	fmt.Println()
	fmt.Printf("Total: %d tarefa(s)\n", len(tasks))
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
