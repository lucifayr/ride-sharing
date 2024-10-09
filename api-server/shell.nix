{ pkgs ? import <nixpkgs> }:

pkgs.mkShell {
  name = "dev-shell";
  buildInputs = with pkgs; [
    sqlc
    sqlite
    go
    nodePackages.sql-formatter
  ];
}
