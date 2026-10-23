{
  description = "Command-line tool for extracting cookies from the user's web browser";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs?tag=24.11";
  };

  outputs = { self, nixpkgs, ... }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in rec {
          cookies = pkgs.buildGoModule rec {
            pname = "cookies";
            version = "0.5.1";
            src = ./.;
            vendorHash = "sha256-lFRZW2KtVsZLHq4oLyDUjFmIpkWuMlIjNuFxLvUp2wg=";
          };

          default = cookies;
        }
      );
    };
}
