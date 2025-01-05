{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages = rec {
          recoredns-ui = with pkgs; stdenv.mkDerivation rec {
            name = "recoredns-ui";
            src = self;
            buildInputs = [
              go
              nodejs
            ];
            buildPhase = ''
              cd web && npm i && npm run build && cd ..
              go get . && go generate ./... && go build . -o recoredns-ui -ldflags "-s -w"
            '';
            installPhase = ''
              mkdir -p $out/bin
              cp recoredns-ui $out/bin
            '';
          };
          default = recoredns-ui;
        };

        devShell = with pkgs; mkShell {
          buildInputs = [
            go
            nodejs
            dig
            tokei
          ];
          GOPATH = "/home/coder/.cache/go";
          RECOREDNS_MYSQL_DSN = "recorednsui:A123456a-@tcp(mysql.dev:3306)/recorednsui?charset=utf8mb4";
        };

        nixosModule = { config, pkgs, lib, ... }: with lib;
          let
            cfg = config.services.recoredns-ui;
          in
          {
            options.services.hangitbot = {
              enable = mkEnableOption "reCoreDNS-UI service";

              mysql-dsn = mkOption {
                type = types.str;
                example = "recorednsui:A123456a-@tcp(mysql.dev:3306)/recorednsui?charset=utf8mb4";
                description = lib.mdDoc "mysql connection DSN";
              };

              extraOptions = mkOption {
                type = types.str;
                description = lib.mdDoc "Extra options";
                default = "";
              };
            };

            config = mkIf cfg.enable {
              systemd.services.recoredns-ui = {
                wantedBy = [ "multi-uesr.target" ];
                environment = {
                  RECOREDNS_MYSQL_DSN = cfg.mysql-dsn;
                };
                serviceconfig.ExecStart = "${pkgs.recoredns-ui}/bin/recoredns-ui server";
              };
            };
          };
      });
}
