package main

import (
	"log"
	"os"

	"github.com/ajr-cabbage/gator/internal/config"
)

func main() {
	// read config file and store config
	c, err := config.Read()
	if err != nil {
		log.Fatalf("%v", err)
	}
	// store config in state{}
	s := &state{conf: c}
	// initialize commands
	cmds := new(commands)
	cmdMap := make(map[string]func(*state, command) error)
	cmds.cmdFuncs = cmdMap
	// register "login" command
	cmds.register("login", handlerLogin)
	// get CLI arguments
	cliArgs := os.Args
	if len(cliArgs) < 2 {
		log.Fatal("No command given.")
	}
	//set cmd name and args
	cmdName := cliArgs[1]
	cmdArgs := cliArgs[2:]
	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
