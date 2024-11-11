{ pkgs ? import <nixpkgs> }:

pkgs.mkShell {
  name = "dev-shell";
  buildInputs = with pkgs; [
    sqlc
    sqlite
    go
    nodePackages.sql-formatter
  ];

  # env vars
  RS_HOST_ADDR = "127.0.0.1:8000";
  RS_WEB_APP_ADDR = "127.0.0.1:5173";
  RS_SECRET_AUTH_TOKEN = "super-secret-fake-token";
  # dev-only client id & secret
  RS_GOOGLE_CLIENT_ID = "750385423567-rsrv4dknuvrts9rv5neab3dl667r5la6.apps.googleusercontent.com";
  RS_GOOGLE_CLIENT_SECRET = "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";
}
