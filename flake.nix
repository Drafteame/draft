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
        draftVersion = "1.8.2";
      in
      {
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            goimports-reviser
            revive
            go-task
          ];
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

          vendorHash = "sha256-GfkC8FL3JsLMsZ8ms9USmvZQE4yJxsr9NoazKO\/5lcE=";

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