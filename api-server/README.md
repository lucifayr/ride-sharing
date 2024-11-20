# ride-sharing API

## dev environment

Using nix to manage dependencies

```sh
nix develop --command $SHELL
```

Or install the following dependencies manually.

- `sqlite`
- `sqlc`
- `go`
- `sql-formatter`

### running (dev)

```sh
go run app/main.go
```

### SQL

To add queries modify/add files in `db/queries/` and run `sqlc generate`.

To update the database schema create a new migration file `db/migrations` by
running

```sh
go run app/cli/main.go migrations create --name {migration-name}
```

Then run `sqlc generate` to update the `golang` types.

#### formatting

```sh
go run app/cli/main.go sql fmt
```
