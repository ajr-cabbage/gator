package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ajr-cabbage/gator/internal/config"
	"github.com/ajr-cabbage/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// read config file and store config
	c, err := config.Read()
	if err != nil {
		log.Fatalf("%v", err)
	}
	// open connection to database
	db, err := sql.Open("postgres", c.DbURL)
	dbQueries := database.New(db)
	// store config in state{}
	s := &state{conf: c, db: dbQueries}
	// initialize commands
	cmds := new(commands)
	cmdMap := make(map[string]func(*state, command) error)
	cmds.cmdFuncs = cmdMap
	// register commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
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
