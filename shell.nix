{ pkgs? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    gcc
    glibc.dev

    # Golang extra
    pkgs.air
    pkgs.templ
    pkgs.sqlc
    pkgs.pgloader
    pkgs.postgresql
    pkgs.sqlitebrowser
  ];

  shellHook = ''
    export CGO_ENABLED=1
    export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
  '';
}
