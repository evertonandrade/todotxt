package main

import (
	"fmt"
	"os"

	"todotxt/internal/cmd"
	"todotxt/internal/store"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Print(cmd.Help())
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	todoFile := os.Getenv("TODO_FILE")
	s := store.New(todoFile)

	var (
		out string
		err error
	)

	switch command {
	case "add", "a":
		out, err = cmd.Add(s, args)
	case "list", "ls":
		out, err = cmd.List(s, args)
	case "do", "x":
		out, err = cmd.Do(s, args)
	case "undo", "unx":
		out, err = cmd.Undo(s, args)
	case "pri", "p":
		out, err = cmd.Priority(s, args)
	case "depri", "dp":
		out, err = cmd.Depri(s, args)
	case "del", "rm":
		out, err = cmd.Delete(s, args)
	case "archive":
		out, err = cmd.Archive(s, args)
	case "help", "-h", "--help":
		fmt.Print(cmd.Help())
		return
	case "version", "-v", "--version":
		fmt.Printf("todotxt %s\n", version)
		return
	default:
		fmt.Fprintf(os.Stderr, "Comando desconhecido: %s\n\n", command)
		fmt.Print(cmd.Help())
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
	if out != "" {
		fmt.Println(out)
	}
}
