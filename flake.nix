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
        draftVersion = "1.7.0";
      in
      {
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

          vendorHash = "sha256-dIZC+UIck1/yc54Qx2ygU8ni+ktiwfkNrCe7ELCJBrU=";

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