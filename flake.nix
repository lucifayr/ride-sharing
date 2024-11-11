{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
      };
      packageLock = builtins.fromJSON (builtins.readFile (./web-app + "/package-lock.json"));
    in
    {
      packages.${system} = {
        api = pkgs.buildGoModule rec {
          name = "api";
          version = "0.0.0";
          src = ./api-server;
          vendorHash = null;
          subPackages = [ "./app/main.go" ];
        };

        web-app = pkgs.buildNpmPackage
          {
            pname = "${packageLock.name}";
            version = "${packageLock.version}";
            src = ./web-app;
            npmDeps = pkgs.importNpmLock {
              npmRoot = ./web-app;
            };
            npmConfigHook = pkgs.importNpmLock.npmConfigHook;
            installPhase = ''
              mkdir $out
              npm run build
              cp -r dist $out/public
            '';
          };
      };
    };
}
