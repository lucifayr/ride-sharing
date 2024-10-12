package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	embeddings "ride_sharing_api"
	utils "ride_sharing_api/app/database"
	"ride_sharing_api/app/database/migrations"

	_ "github.com/mattn/go-sqlite3"
)

type Command struct {
	subcommands map[string]Command
	exec        func(cmd Command, args []string)
}

var commands = Command{
	subcommands: map[string]Command{
		"migrations": {
			subcommands: map[string]Command{
				"up": {
					exec: func(_cmd Command, _args []string) {
						runMigrationsUp()
					},
				},
				"down": {
					exec: func(_cmd Command, _args []string) {
						runMigrationsDown()
					},
				},
			},
			exec: func(cmd Command, _args []string) {
				unfinishedCmd(cmd)
			},
		},
		"sql": {
			subcommands: map[string]Command{
				"fmt": {
					exec: func(_cmd Command, _args []string) {
						fmtSqlFiles()
					},
				},
			},
			exec: func(cmd Command, _args []string) {
				unfinishedCmd(cmd)
			},
		},
	},
	exec: func(cmd Command, _args []string) {
		unfinishedCmd(cmd)
	},
}

func main() {
	args := os.Args[1:]
	parseCmd(commands, args)
}

func runMigrationsUp() {
	db := setupDb()
	m := migrations.FromEmbedFs(embeddings.DbMigrations, "db/migrations")
	m.Up(db)
}

func runMigrationsDown() {
	db := setupDb()
	m := migrations.FromEmbedFs(embeddings.DbMigrations, "db/migrations")
	m.Down(db)
}

func fmtSqlFiles() {
	filepath.Walk("db", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatalln("Failed to read file.", "path", path, "error", err)
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".sql") {
			return nil
		}

		fmt.Println("Formatting SQL file :", path)
		fmtCmd := exec.Command("sql-formatter", "--config", ".sql-formatter.json", "--fix", path)
		err = fmtCmd.Run()
		if err != nil {
			log.Fatalln("Failed to format SQL file with 'sql-formatter'.", "path", path, "error", err)
		}

		return nil
	})
}

func setupDb() *sql.DB {
	dbFile := "./database.db"
	err := utils.CreateDbFileIfNotExists(dbFile)
	if err != nil {
		log.Fatalln("Failed to create database file.", dbFile, err)
	}

	db, err := sql.Open("sqlite3", "file:"+dbFile)
	if err != nil {
		log.Fatalln("Failed to connect to database.", dbFile, err)
	}

	return db
}

func parseCmd(cmd Command, args []string) {
	if len(args) == 0 {
		cmd.exec(cmd, args)
		return
	}

	sub, ok := cmd.subcommands[args[0]]
	if !ok {
		cmd.exec(cmd, args)
		return
	}

	parseCmd(sub, args[1:])
}

func unfinishedCmd(cmd Command) {
	cmds := make([]string, len(cmd.subcommands))

	i := 0
	for k := range cmd.subcommands {
		cmds[i] = k
		i++
	}

	fmt.Printf("Missing sub-command argument. Possible options are [%s].\n", strings.Join(cmds, ", "))
	os.Exit(1)
}
