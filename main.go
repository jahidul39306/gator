package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jahidul39306/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandNames map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.commandNames[cmd.name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.name)
}

func (c *commands) register(name string, handler func(*state, command) error) {
	c.commandNames[name] = handler
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("username is required")
	}
	username := cmd.arguments[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User '%s' has been set\n", username)
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	sta := state{cfg: cfg}
	cmds := commands{commandNames: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal(fmt.Errorf("not enough arguments were provided"))
	}
	cmd := command{name: args[0], arguments: args[1:]}
	err = cmds.run(&sta, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
