{pkgs}:
pkgs.buildGoModule {
  name = "freeze";
  src = ./.;
  vendorHash = "sha256-6++AGNKse/cBg5TAybnoiuvgMf8KjXJ37fVwFNDK4Ic=";
}
