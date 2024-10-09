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


### running

```sh
go run ride-sharing-api/main.go
```

### SQL

- To add queries modify files in `db/queries/{table-name}.sql` and run `sqlc
  generate`.
- To update the database schema modify files in
  `db/migrations/{table-name}.[up/down].sql` and run `sqlc generate`.

#### formatting

```sh
sh ./sql-fmt.sh
```
