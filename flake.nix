{
  description = "Resec flake";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  outputs = { self, nixpkgs }: let
    inherit (nixpkgs.legacyPackages.x86_64-linux) lib buildGoModule;
  in {
    packages.x86_64-linux.resec = buildGoModule rec {
      pname = "resec";
      version = "master-${lib.substring 0 7 self.rev}";

      src = ./.;
      vendorSha256 = "sha256-r1RIYpX9akvoxW8qav+mJbpRJwSId4Pg0xE4kbSOtFA=";

      # Checks freeze
      doCheck = false;

      ldflags = [
        "-s" "-w"
        "-X main.Version=master-${self.rev}"
      ];

      meta = with lib; {
        description = "ReSeC- Redis Service Consul";
        homepage = "https://github.com/johnalotoski/resec";
        platforms = platforms.linux;
        license = licenses.mit;
        maintainers = with maintainers; [ ];
      };
    };
  };
}
