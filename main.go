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
		cmd.Help()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	todoFile := os.Getenv("TODO_FILE")
	s := store.New(todoFile)

	switch command {
	case "add", "a":
		cmd.Add(s, args)
	case "list", "ls":
		cmd.List(s, args)
	case "do", "x":
		cmd.Do(s, args)
	case "undo", "unx":
		cmd.Undo(s, args)
	case "pri", "p":
		cmd.Priority(s, args)
	case "depri", "dp":
		cmd.Depri(s, args)
	case "del", "rm":
		cmd.Delete(s, args)
	case "archive":
		cmd.Archive(s, args)
	case "help", "-h", "--help":
		cmd.Help()
	case "version", "-v", "--version":
		fmt.Printf("todotxt %s\n", version)
	default:
		fmt.Fprintf(os.Stderr, "Comando desconhecido: %s\n\n", command)
		cmd.Help()
		os.Exit(1)
	}
}
