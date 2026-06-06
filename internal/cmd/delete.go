package cmd

import (
	"fmt"
	"strconv"

	"todotxt/internal/store"
)

func Delete(s *store.Store, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("forneça o número da tarefa. Uso: todotxt del <número> [número...]")
	}

	tasks, err := loadTasks(s)
	if err != nil {
		return "", fmt.Errorf("ao carregar tarefas: %w", err)
	}
	if len(tasks) == 0 {
		return "", fmt.Errorf("nenhuma tarefa para remover")
	}

	indices := map[int]bool{}
	for _, a := range args {
		n, err := strconv.Atoi(a)
		if err != nil || n < 1 || n > len(tasks) {
			return "", fmt.Errorf("número inválido %q", a)
		}
		indices[n-1] = true
	}

	removed := 0
	filtered := tasks[:0]
	for i, t := range tasks {
		if indices[i] {
			removed++
			continue
		}
		filtered = append(filtered, t)
	}
	tasks = filtered

	if removed == 0 {
		return "Nenhuma tarefa removida.", nil
	}

	if err := saveTasks(s, tasks); err != nil {
		return "", fmt.Errorf("ao guardar tarefas: %w", err)
	}
	return fmt.Sprintf("%d tarefa(s) removida(s).", removed), nil
}
