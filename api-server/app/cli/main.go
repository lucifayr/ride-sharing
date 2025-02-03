package main

import (
	"context"
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
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/common"
	"ride_sharing_api/app/database/migrations"
	"ride_sharing_api/app/rest"
	"ride_sharing_api/app/sqlc"
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
		"dev": {
			subcommands: map[string]Command{
				"make-account": {
					exec: func(_cmd Command, args []string) {
						id := parseCmdFlag(args, "--id", "-i")
						email := parseCmdFlag(args, "--email", "-e")
						name := parseCmdFlag(args, "--name", "-n")

						dbFile := utils.GetEnvRequired(common.ENV_DB_NAME)
						db, err := utils.InitDb(dbFile)
						assert.Nil(err)

						tx, err := db.Begin()
						assert.Nil(err)

						queries := sqlc.New(tx)
						_, err = queries.UsersCreate(context.Background(), sqlc.UsersCreateParams{ID: id, Name: name, Email: email, Provider: "google"})
						assert.Nil(err)

						tokens := rest.GenAuthTokens(id, email)

						params := sqlc.UsersSetTokensParams{ID: id, AccessToken: utils.SqlNullStr(tokens.AccessToken), RefreshToken: utils.SqlNullStr(tokens.RefreshToken)}
						err = queries.UsersSetTokens(context.Background(), params)
						assert.Nil(err)

						err = tx.Commit()
						assert.Nil(err)

						log.Printf(`
id: "%s",
name: "%s",
email: "%s",
tokens: {
	accessToken: "%s",
	refreshToken: "%s",
}
`, id, name, email, tokens.AccessToken, tokens.RefreshToken)
					},
				},
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

		content := []byte{}
		if strings.Contains(suffix, migrations.MigrationKindValidate) {
			content = []byte("SELECT * FROM missing_validation_migration LIMIT 1;")
		}

		err := os.WriteFile(path, content, 0644)
		if err != nil {
			log.Fatalln("Failed to write migrations file.", "path", path, "error", err)
		}

		fmt.Println("-", path)
	}
}

func parseCreateMigrationFlags(args []string) map[string]any {
	name := parseCmdFlag(args, "--name", "-n")

	return map[string]any{
		"name": name,
	}
}

func parseCmdFlag(args []string, long string, short string) string {
	idx := utils.IdxOf(args, func(arg string) bool {
		return arg == long || arg == short
	})

	if idx == -1 {
		log.Fatalf("Missing required parameter	`%s`. Provided it with `%s <value>`.\n", long, long)
	}

	value, ok := utils.SliceGet(args, idx+1)
	if !ok {
		log.Fatalf("Missing required parameter	`%s`. Provided it with `%s <value>`.\n", long, long)
	}

	return *value
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
