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
    in
    {
      packages.${system} = {
        web-app = pkgs.buildNpmPackage
          {
            name = "web-app";
            src = self;
            installPhase = ''
              mkdir $out
              npm run build
              cp -r dist $out/public
            '';
          };
      };
    };
}
