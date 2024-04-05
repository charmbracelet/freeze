{
  description = "A tool for generating images of code and terminal output";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {inherit system;};
    in {
      packages.default = import ./default.nix {inherit pkgs;};
    })
    // {
      overlays.default = final: prev: {
        freeze = import ./default.nix {pkgs = final;};
      };
    };
}
