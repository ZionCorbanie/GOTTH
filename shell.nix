{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gnumake
    nodejs_26
    cargo
  ];

  shellHook = ''
  '';
}

