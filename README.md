# reCoreDNS-UI

Web UI for CoreDNS
Initial code from [reCoreD-UI](https://github.com/reCoreD-UI/reCoreD-UI) with some bug fixes and small improvements.

## UI

![ui](.assets/ui.webp)

## Build locally

Install `go` and `npm` first.

```bash
# Build web first
(cd web && npm run build)

# Build server
go get .
go generate ./...
go build .
```

## Running locally

Build [coredns](https://coredns.io/) with [mysql](coredns.io/explugins/mysql/) plugin first.

A mysql server is needed.

```bash
# example
export RECOREDNS_MYSQL_DSN="recorednsui:A123456a-@tcp(mysql.dev:3306)/recorednsui?charset=utf8mb4"
export RECOREDNS_MYSQL_DSN="coredns:coredns@tcp(mysql.fontanas-uk.tailcloud:3306)/coredns?charset=utf8mb4"
./reCoreDNS-UI config db migrate

# setup admin user
./reCoreDNS-UI config user -u user -p password

# setup DNS
./reCoreDNS-UI config dns -s 1.1.1.1 -s 1.2.3.4

# run server and open http://localhost:3000
./reCoreDNS-UI server
```

```ini
# systemd service
[Unit]
Description=reCoreDNS-UI

[Service]
Type=simple
# RECOREDNS_MYSQL_DSN="dsn"
EnvironmentFile=-/etc/default/recoredns-ui
EnvironmentFile=-/etc/sysconfig/recoredns-ui
ExecStart=/usr/local/bin/reCoreDNS-UI server

[Install]
WantedBy=multi-user.target
```

## Docker image
Install `docker` first.

`docker build -t reCoreDNS-UI:latest .`

Once the build is succesfull, connect to the container 
`docker run -d --name reCoreDNS-UI -it salma:latest /bin/sh`

and run the following commands:
- Preparet eh database
    `./reCoreDNS-UI config db migrate`
- Configure login credentials to access reCoreDNS-UI
    `./reCoreDNS-UI config user -u $USER -p $PASSWORD`
- Setup the DNS:
    `./reCoreDNS-UI config dns -s 1.1.1.1 -s 1.2.3.4`

Stop the container without deleting it (it's sufficient to type `exit` after running the last command) then restart it using

`docker start reCoreDNS-UI`

Head over to http://localhost:3000 and to access reCoreDNS-UI and start managing your DNS entries.

## Docker Compose