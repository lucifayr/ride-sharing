{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkg-version = "0.0.0";
      pkgs = import nixpkgs {
        inherit system;
      };
    in
    {
      packages.${system} = {
        default =
          import ./shell.nix { inherit pkgs; };

        api = pkgs.buildGoModule rec {
          name = "api";
          version = pkg-version;
          src = ./.;
          vendorHash = null;
          subPackages = [ "./app/main.go" ];
        };
      };
    };
}
