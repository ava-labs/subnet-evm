{
  # To use:
  #  - install nix: https://github.com/DeterminateSystems/nix-installer?tab=readme-ov-file#install-nix
  #  - run `nix develop` or use direnv (https://direnv.net/)
  #    - for quieter direnv output, set `export DIRENV_LOG_FORMAT=`

  description = "Subnet-EVM development environment";

  inputs = {
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.2405.*.tar.gz";
    avalanchego.url = "github:ava-labs/avalanchego?ref=198b68f0a850fbfa12e50735bed56b14d99fe0f1";
  };

  outputs = { self, nixpkgs, avalanchego, ... }:
    let
      allSystems = builtins.attrNames avalanchego.devShells;
      forAllSystems = nixpkgs.lib.genAttrs allSystems;
    in {
      # Define the development shells for this repository
      devShells = forAllSystems (system: {
        default = avalanchego.devShells.${system}.default;
      });
    };
}
