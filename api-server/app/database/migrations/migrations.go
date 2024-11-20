package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path"
	"ride_sharing_api/app/assert"
	"slices"
	"sort"
	"strings"
)

const (
	MigrationKindUp       = "up"
	MigrationKindValidate = "validate"
)

type Migrations struct {
	migrationSourceTexts []Migration
}

type Migration struct {
	name string
	sql  map[string]string
}

func FileSuffixes() []string {
	return []string{
		fmt.Sprintf("%s.sql", MigrationKindUp),
		fmt.Sprintf("%s.sql", MigrationKindValidate),
	}
}

func FromEmbedFs(fs embed.FS, root string) Migrations {
	dir, err := fs.ReadDir(root) // not using simulator because embed.fs is not doing any IO (basically)
	if err != nil {
		log.Fatalln("Failed to read migrations from root directory.", "root", root)
	}

	var migrations []Migration

	for _, entry := range dir {
		info, err := entry.Info()
		assert.True(err == nil, "Failed to get file system info of migration entry.", "root", root, "entry", entry.Name())
		assert.True(!info.IsDir(), "Expected all entries in migration directory to be sql files. Found sub-directory.", "root", root, "sub-dir", info.Name())

		name := info.Name()

		mUpSuffix := fmt.Sprintf(".%s.sql", MigrationKindUp)
		isUp := strings.HasSuffix(name, mUpSuffix)

		mValSuffix := fmt.Sprintf(".%s.sql", MigrationKindValidate)
		isVal := strings.HasSuffix(name, mValSuffix)

		assert.True(isUp || isVal, "Expected file to end with '.up.sql' or '.validate.sql'.", "name", name)

		path := path.Join(root, name)

		file, err := fs.ReadFile(path) // not using simulator because of embedded fs
		assert.True(err == nil, "Failed to read migrations file.", "root", root, "path", path)

		sqlSource := string(file[:])

		var mName string
		var mKind string
		if isUp {
			mName = strings.Replace(name, mUpSuffix, "", 1)
			mKind = MigrationKindUp
		}

		if isVal {
			mName = strings.Replace(name, mValSuffix, "", 1)
			mKind = MigrationKindValidate
		}

		assert.False(strings.HasSuffix(mName, ".sql"), "Failed to strip prefix from migrations file.", "path", path, "migration name", mName)

		idx := slices.IndexFunc(migrations, func(m Migration) bool {
			return m.name == mName
		})

		var m Migration
		if idx == -1 {
			sql := make(map[string]string)
			m = Migration{name: mName, sql: sql}
			migrations = append(migrations, m)
		} else {
			m = migrations[idx]
		}

		m.sql[mKind] = sqlSource
	}

	for _, m := range migrations {
		assert.True(contains(m.sql, MigrationKindUp), "Migration is missing .up variant", "name", m.name)
		assert.True(contains(m.sql, MigrationKindValidate), "Migration is missing .validate variant", "name", m.name)
	}

	sort.Slice(migrations, func(i, j int) bool {
		a := migrations[i]
		b := migrations[j]

		return a.name < b.name
	})

	return Migrations{migrationSourceTexts: migrations}
}

func (migrations *Migrations) Up(db *sql.DB) {
	for _, mig := range migrations.migrationSourceTexts {
		sqlVal := mig.sql[MigrationKindValidate]
		_, err := db.Exec(sqlVal)
		if err == nil {
			log.Printf("Skipping migration, schema passed validation [%s] : %s\n", MigrationKindUp, mig.name)
			continue
		}

		sql := mig.sql[MigrationKindUp]

		log.Printf("Running migration [%s] : %s\n", MigrationKindUp, mig.name)
		err = runMigration(db, sql)
		assert.Nil(err, "Migration failed.", "name", mig.name)
	}
}

func runMigration(db *sql.DB, sql string) error {
	res, err := db.Exec(sql)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Println("\t- Failed to get affected Rows", "error", err)
	} else {
		log.Println("\t- Rows affected", affected)
	}

	return nil
}

func contains(m map[string]string, key string) bool {
	_, ok := m[key]
	return ok
}
