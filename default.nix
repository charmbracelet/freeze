{pkgs}:
pkgs.buildGoModule {
  name = "freeze";
  src = ./.;
  vendorHash = "sha256-OFNpZ6BOxC1nVmf0X89wlzIb3S7xS89YX4EkYXF4ozM=";
}
