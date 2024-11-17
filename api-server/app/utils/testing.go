package utils

import (
	"context"
	"os"
	"path"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/common"
	"ride_sharing_api/app/sqlc"
	"sync"
	"syscall"
)

var once sync.Once

func SetupTestDBs() {
	once.Do(func() {
		regenerateTestDbs()
	})
}

func regenerateTestDbs() {
	dbDir := path.Join(ProjectRoot(), "db/testing")

	syscall.Umask(0)
	err := os.RemoveAll(dbDir)
	assert.Nil(err)
	err = os.Mkdir(dbDir, 0777)
	assert.Nil(err)

	for db, users := range common.TEST_USERS {
		dbName := path.Join(dbDir, db+".db")
		err = CreateDbFileIfNotExists(dbName)
		assert.Nil(err)

		db, err := InitDb(dbName)
		assert.Nil(err)
		queries := sqlc.New(db)

		for _, user := range users {
			args_create := sqlc.UsersCreateParams{
				ID:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Provider: user.Provider,
			}
			_, err = queries.UsersCreate(context.Background(), args_create)
			assert.Nil(err)

			args_set_tokens := sqlc.UsersSetTokensParams{
				AccessToken:  user.AccessToken,
				RefreshToken: user.RefreshToken,
				ID:           user.ID,
			}
			err = queries.UsersSetTokens(context.Background(), args_set_tokens)
			assert.Nil(err)
		}

	}
}
