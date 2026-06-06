package cmd

import (
	"fmt"

	"todotxt/internal/store"
	"todotxt/internal/task"
)

func Do(s *store.Store, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("forneça o número da tarefa. Uso: todotxt do <número>")
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
		return "", fmt.Errorf("tarefa %d não existe (existem %d)", num, len(tasks))
	}

	t := tasks[num-1]
	if t.Completed {
		return "", fmt.Errorf("tarefa %d já está concluída", num)
	}
	t.MarkComplete()
	if err := saveTasks(s, tasks); err != nil {
		return "", fmt.Errorf("ao guardar tarefas: %w", err)
	}
	return fmt.Sprintf("Tarefa %d concluída: %s", num, t.Description), nil
}

func Undo(s *store.Store, args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("forneça o número da tarefa. Uso: todotxt undo <número>")
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
	if !t.Completed {
		return "", fmt.Errorf("tarefa %d não está concluída", num)
	}
	t.MarkIncomplete()
	if err := saveTasks(s, tasks); err != nil {
		return "", fmt.Errorf("ao guardar tarefas: %w", err)
	}
	return fmt.Sprintf("Tarefa %d reaberta: %s", num, t.Description), nil
}
