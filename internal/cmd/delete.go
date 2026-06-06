package cmd

import (
	"fmt"
	"os"
	"strconv"

	"todotxt/internal/store"
)

func Delete(s *store.Store, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Erro: forneça o número da tarefa.")
		fmt.Fprintln(os.Stderr, "Uso: todotxt del <número> [número...]")
		os.Exit(1)
	}

	tasks := loadTasks(s)
	if len(tasks) == 0 {
		fmt.Fprintln(os.Stderr, "Nenhuma tarefa para remover.")
		os.Exit(1)
	}

	indices := map[int]bool{}
	for _, a := range args {
		n, err := strconv.Atoi(a)
		if err != nil || n < 1 || n > len(tasks) {
			fmt.Fprintf(os.Stderr, "Erro: número inválido %q\n", a)
			os.Exit(1)
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
		fmt.Println("Nenhuma tarefa removida.")
		return
	}

	saveTasks(s, tasks)
	printOK(fmt.Sprintf("%d tarefa(s) removida(s).", removed))
}
