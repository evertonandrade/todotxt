package cmd

import (
	"fmt"
	"os"

	"todotxt/internal/store"
	"todotxt/internal/task"
)

func Priority(s *store.Store, args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Erro: forneça o número da tarefa e a prioridade.")
		fmt.Fprintln(os.Stderr, "Uso: todotxt pri <número> <A-Z>")
		os.Exit(1)
	}
	num, err := task.ParseLineNumber(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
	pri, err := task.ParsePriority(args[1])
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
	old := t.Priority
	t.SetPriority(pri)
	saveTasks(s, tasks)

	if old == "" && pri == "" {
		fmt.Println("Nenhuma alteração.")
	} else if pri == "" {
		printOK(fmt.Sprintf("Prioridade da tarefa %d removida.", num))
	} else {
		printOK(fmt.Sprintf("Prioridade da tarefa %d alterada: %s → %s", num, old, pri))
	}
}

func Depri(s *store.Store, args []string) {
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
	if t.Priority == "" {
		fmt.Println("Tarefa sem prioridade.")
		return
	}
	t.SetPriority("")
	saveTasks(s, tasks)
	printOK(fmt.Sprintf("Prioridade da tarefa %d removida.", num))
}
