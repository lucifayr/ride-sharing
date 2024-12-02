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
  RS_DB_NAME = "rides.db";
  RS_HOST_ADDR = "127.0.0.1:8000";
  RS_CORS_ORIGIN = "http://127.0.0.1";
  RS_WEB_APP_URL = "http://127.0.0.1:5173";
  RS_SECRET_AUTH_TOKEN = "03**CsL@pfFmtt5K4LE*SVYXPseFZ^FO";
  RS_NO_TLS = "true";
  # dev-only client id & secret
  RS_GOOGLE_REDIRECT_URL = "http://127.0.0.1:8000/auth/google/callback";
  RS_GOOGLE_CLIENT_ID = "750385423567-8vu2cst8njm4d6ple8e424ltpd9dh9t2.apps.googleusercontent.com";
  RS_GOOGLE_CLIENT_SECRET = "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";
}
