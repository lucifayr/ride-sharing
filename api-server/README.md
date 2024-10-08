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


### running

```sh
go run ride-sharing-api/main.go
```

### SQL

To add queries modify `query.sql` and run `sqlc generate`.
