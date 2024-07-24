{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: {
  env.CGO_ENABLED = true;

  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    gopls
    gdlv
    delve
    golangci-lint
    libcap
    gcc
    go
  ];

  # https://devenv.sh/languages/
  languages = {
    #nix.enable = true;
    go = {
      enable = true;
      # package = pkgs.go_1_23;
    };
  };

  # https://devenv.sh/pre-commit-hooks/
  pre-commit.hooks = {
    shellcheck.enable = true;
  };
}
