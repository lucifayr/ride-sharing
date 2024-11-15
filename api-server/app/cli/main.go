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
	"time"

	embeddings "ride_sharing_api"
	"ride_sharing_api/app/common"
	"ride_sharing_api/app/database/migrations"
	utils "ride_sharing_api/app/utils"

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
				"create": {
					exec: func(_cmd Command, args []string) {
						flags := parseCreateMigrationFlags(args)
						createMigration(flags)
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

func createMigration(flags map[string]any) {
	prefix := time.Now().Local().Format(time.DateTime)
	// Replace space in date time string of format "2006-01-02 15:04:05"
	prefix = strings.Replace(prefix, " ", "_", 1)
	// Replace ':' in date time string of format "2006-01-02 15:04:05"
	prefix = strings.Replace(prefix, ":", "-", 2)

	migName := *flags["name"].(*string)
	name := fmt.Sprintf("%s_%s", prefix, migName)

	for _, suffix := range migrations.FileSuffixes() {
		path := fmt.Sprintf("./db/migrations/%s.%s", name, suffix)
		if utils.FileExists(path) {
			log.Fatalln("Migration already exists", "path", path)
		}

		err := os.WriteFile(path, []byte{}, 0644)
		if err != nil {
			log.Fatalln("Failed to write migrations file.", "path", path, "error", err)
		}

		fmt.Println("-", path)
	}
}

func parseCreateMigrationFlags(args []string) map[string]any {
	idx := utils.IdxOf(args, func(arg string) bool {
		return arg == "--name" || arg == "-n"
	})

	if idx == -1 {
		log.Fatalln("Missing required parameter	`--name`. Provided it with `--name <value>`.")
	}

	name, ok := utils.SliceGet(args, idx+1)
	if !ok {
		log.Fatalln("Missing value for required parameter `--name`. Provided it with `--name <value>`.")
	}

	return map[string]any{
		"name": name,
	}
}

func setupDb() *sql.DB {
	dbFile := utils.GetEnvRequired(common.ENV_DB_NAME)
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
