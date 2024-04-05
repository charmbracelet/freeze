{pkgs}:
pkgs.buildGoModule {
  name = "freeze";
  src = ./.;
  vendorHash = "sha256-AUFzxmQOb/h0UgcprY09IVI7Auitn3JTDU/ptKicIAU=";
}
