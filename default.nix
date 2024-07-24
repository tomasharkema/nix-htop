let
  pkgs = import <nixpkgs> {};
in {
  nix-htop = pkgs.callPackage ./package.nix {};
}
