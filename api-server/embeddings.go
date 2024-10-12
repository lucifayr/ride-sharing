package embeddings

import "embed"

//go:embed db/migrations/*
var DbMigrations embed.FS
