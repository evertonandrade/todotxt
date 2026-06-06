package cmd

import (
	"fmt"
	"os"
	"strings"

	"todotxt/internal/store"
)

func Archive(s *store.Store, _ []string) {
	count, err := s.Archive()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao arquivar: %v\n", err)
		os.Exit(1)
	}
	if count == 0 {
		fmt.Println("Nenhuma tarefa concluída para arquivar.")
		return
	}
	printOK(fmt.Sprintf("%d tarefa(s) movida(s) para %s.", count, s.DoneFile))
}

func Help() {
	help := `todotxt — gestor de tarefas CLI compatível com todo.txt

Uso:
  todotxt <comando> [argumentos]

Comandos:
  add "descrição" [+proj] [@ctx] [pri:A] [due:YYYY-MM-DD]
      Adiciona uma nova tarefa.

  list [filtros]
      Lista tarefas pendentes. Filtros:
        +projeto      filtra por projeto
        @contexto     filtra por contexto
        pri:A         filtra por prioridade (A-Z)
        all           mostra pendentes e concluídas
        done          mostra apenas concluídas
        due           mostra apenas com data de vencimento
        overdue       mostra apenas vencidas
        today         mostra apenas que vencem hoje
        texto         pesquisa textual na descrição

  do <número>        marca tarefa como concluída
  undo <número>      reabre uma tarefa concluída

  pri <número> <A-Z> define prioridade (use "-" para remover)
  depri <número>     remove a prioridade

  del <número>...    remove uma ou mais tarefas
  archive            move concluídas para done.txt

  help               mostra esta ajuda
  version            mostra a versão

Ficheiros:
  TODO_DIR/todo.txt  por predefinição ./todo.txt
  Arquivo:           por predefinição ./done.txt (ou TODO_DIR/done.txt)

Variáveis de ambiente:
  TODO_DIR           diretório dos ficheiros todo.txt e done.txt
  NO_COLOR           desativa saída colorida
`
	fmt.Print(strings.TrimLeft(help, "\n"))
}
