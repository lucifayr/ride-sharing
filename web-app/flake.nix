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
      packageLock = builtins.fromJSON (builtins.readFile (./. + "/package-lock.json"));
    in
    {
      packages.${system} = {
        web-app = pkgs.buildNpmPackage
          {
            pname = "${packageLock.name}";
            version = "${packageLock.version}";
            src = ./.;
            npmDeps = pkgs.importNpmLock {
              npmRoot = ./.;
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
