{
  buildGoModule,
  musl,
}:
buildGoModule rec {
  pname = "nix-htop";
  version = "0.0.1";

  src =
    builtins.filterSource
    (path: type: !(type == "directory" && (baseNameOf path == ".devenv" || baseNameOf path == ".direnv")))
    ./.;

  vendorHash = "sha256-XX0AMpZR3v9bZdLiqDYVT0pOGBbrM7rJERTmxlFIgZo=";

  CGO_ENABLED = 1;
}
