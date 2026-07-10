package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Expected 'serve' or 'migrate' subcommand")
	}

	switch os.Args[1] {
	case "serve":
		log.Println("Launching HTTP Server...")
		err := runCommand("go", "run", "./cmd/server/main.go")
		if err != nil {
			panic(err)
		}
	case "migrate":
		log.Println("Running Database Migrations...")
		err := runCommand("go", "run", "./cmd/migration/main.go")
		if err != nil {
			panic(err)
		}
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func runCommand(command string, args ...string) error {
    cmd := exec.Command(command, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}