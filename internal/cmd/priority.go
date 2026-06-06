package cmd

import (
	"fmt"

	"todotxt/internal/store"
	"todotxt/internal/task"
)

func Priority(s *store.Store, args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("forneça o número da tarefa e a prioridade. Uso: todotxt pri <número> <A-Z>")
	}
	num, err := task.ParseLineNumber(args[0])
	if err != nil {
		return "", err
	}
	pri, err := task.ParsePriority(args[1])
	if err != nil {
		return "", err
	}

	tasks, err := loadTasks(s)
	if err != nil {
		return "", fmt.Errorf("ao carregar tarefas: %w", err)
	}
	if num < 1 || num > len(tasks) {
		return "", fmt.Errorf("tarefa %d não existe", num)
	}
	t := tasks[num-1]
	old := t.Priority
	t.SetPriority(pri)
	if err := saveTasks(s, tasks); err != nil {
		return "", fmt.Errorf("ao guardar tarefas: %w", err)
	}

	if old == "" && pri == "" {
		return "Nenhuma alteração.", nil
	}
	if pri == "" {
		return fmt.Sprintf("Prioridade da tarefa %d removida.", num), nil
	}
	return fmt.Sprintf("Prioridade da tarefa %d alterada: %s → %s", num, old, pri), nil
}

func Depri(s *store.Store, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("forneça o número da tarefa. Uso: todotxt depri <número>")
	}
	num, err := task.ParseLineNumber(args[0])
	if err != nil {
		return "", err
	}

	tasks, err := loadTasks(s)
	if err != nil {
		return "", fmt.Errorf("ao carregar tarefas: %w", err)
	}
	if num < 1 || num > len(tasks) {
		return "", fmt.Errorf("tarefa %d não existe", num)
	}
	t := tasks[num-1]
	if t.Priority == "" {
		return "Tarefa sem prioridade.", nil
	}
	t.SetPriority("")
	if err := saveTasks(s, tasks); err != nil {
		return "", fmt.Errorf("ao guardar tarefas: %w", err)
	}
	return fmt.Sprintf("Prioridade da tarefa %d removida.", num), nil
}
