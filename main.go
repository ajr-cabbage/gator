package main

import (
	"database/sql"
	"log"

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
	cmds := getCommands()
	// build command from os.Args
	cmd := buildCommand()
	// run retrieved command
	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
