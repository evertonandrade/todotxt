package cmd

import (
	"fmt"
	"os"

	"todotxt/internal/store"
	"todotxt/internal/task"
)

func Do(s *store.Store, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Erro: forneça o número da tarefa.")
		fmt.Fprintln(os.Stderr, "Uso: todotxt do <número>")
		os.Exit(1)
	}
	num, err := task.ParseLineNumber(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	tasks := loadTasks(s)
	if num < 1 || num > len(tasks) {
		fmt.Fprintf(os.Stderr, "Erro: tarefa %d não existe (existem %d).\n", num, len(tasks))
		os.Exit(1)
	}

	t := tasks[num-1]
	if t.Completed {
		fmt.Fprintf(os.Stderr, "Tarefa %d já está concluída.\n", num)
		return
	}
	t.MarkComplete()
	saveTasks(s, tasks)
	printOK(fmt.Sprintf("Tarefa %d concluída: %s", num, t.Description))
}

func Undo(s *store.Store, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Erro: forneça o número da tarefa.")
		os.Exit(1)
	}
	num, err := task.ParseLineNumber(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}

	tasks := loadTasks(s)
	if num < 1 || num > len(tasks) {
		fmt.Fprintf(os.Stderr, "Erro: tarefa %d não existe.\n", num)
		os.Exit(1)
	}
	t := tasks[num-1]
	if !t.Completed {
		fmt.Fprintf(os.Stderr, "Tarefa %d não está concluída.\n", num)
		return
	}
	t.MarkIncomplete()
	saveTasks(s, tasks)
	printOK(fmt.Sprintf("Tarefa %d reaberta: %s", num, t.Description))
}
