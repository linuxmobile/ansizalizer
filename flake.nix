{
  description = "A TUI to convert Images to ANSI strings";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages = rec {
        default = ansizalizer;

        ansizalizer = pkgs.buildGoModule {
          pname = "ansizalizer";
          version = "0.1.0";

          src = ./.;

          vendorHash = "sha256-rjwb8+AChOZx3UCAukW3VS7noDq6Jrgob6gySHxmXJI=";

          tags = ["linux"];

          preBuild = ''
            rm env/os_unix.go
          '';

          env = {
            CGO_ENABLED = "1";
          };

          buildInputs = with pkgs; [
            pkg-config
          ];

          meta = with pkgs.lib; {
            description = "A TUI to convert Images to ANSI strings using bubbletea";
            homepage = "https://github.com/Zebbeni/ansizalizer";
            license = licenses.mit;
            maintainers = with maintainers; [linuxmobile];
            mainProgram = "ansizalizer";
          };
        };
      };

      apps = rec {
        default = ansizalizer;
        ansizalizer = flake-utils.lib.mkApp {
          drv = self.packages.${system}.ansizalizer;
        };
      };
    });
}
