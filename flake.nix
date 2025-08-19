{
  description = "Draft CLI tool";

  inputs = {
      nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
      flake-utils.url = "github:numtide/flake-utils";
    };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        draftVersion = "1.13.0";
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            goimports-reviser
            revive
            go-task
          ];

          shellHook = ''
            export GOROOT="${pkgs.go}/share/go"
          '';
        };

        packages.default = pkgs.buildGoModule {
          pname = "draft";
          version = draftVersion;

          subPackages = [ "." ];

          src = ./.;

          ldflags = [
            "-s"
            "-w"
            "-X 'github.com/Drafteame/draft/cmd/commands.Version=nix-v${draftVersion}'"
          ];

          env.CGO_ENABLED = false;
          env.GOWORK = "off";

          vendorHash = "sha256-di5PXoshB1bBmo+qas6Auzh8ma7BWkgyos0aCJ3Q\/k0=";

          meta = {
            description = "CLI tool for creating services and lambdas on Draftea monorepo";
            mainProgram = "draft";
          };
        };

        apps.default = flake-utils.lib.mkApp {
          drv = self.packages.${system}.default;
        };
      }
    );
}