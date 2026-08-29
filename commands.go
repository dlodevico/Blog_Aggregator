package main

import (
	"fmt"
)

// command represents a CLI command name and its parameters.
type command struct {
	Name string
	Args []string
}

// commands maintains the map of supported commands to their handlers.
type commands struct {
	registeredCmds map[string]func(*state, command) error
}

// register binds a command name to a handler function.
func (c *commands) register(name string, f func(*state, command) error) {
	if c.registeredCmds == nil {
		c.registeredCmds = make(map[string]func(*state, command) error)
	}
	c.registeredCmds[name] = f
}

// run executes the handler registered to the provided command name.
func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registeredCmds[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}
