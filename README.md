xbsapi
------

An alternative REST API that serves requests from 
[xBrowserSync](https://www.xbrowsersync.org/) client apps, that is a single 
binary and supports SQLite3, PostgreSQL and MySQL.


## Build

```sh
$ make install-deps
$ make
```


## Run

```sh
$ xbsapi
```


## Configuration

xbsapi will read its config either from a file or from environment
variables. Every configuration key available in the 
[`lib/config.go`](lib/config.go) can be exported as
environment variable, by separating scopes using `_` and prepend `XBSAPI` to
it. For example, the following configuration:

```toml
[Server]
BindIP = "0.0.0.0"
```

... can also be specified as an environment variable:

```sh
export XBSAPI_SERVER_BINDIP="0.0.0.0"
```

xbsapi will try to read the `xbsapi.toml` file from one of the following
paths:

- `/etc/xbsapi.toml`
- `$XDG_CONFIG_HOME/xbsapi.toml`
- `$HOME/.config/xbsapi.toml`
- `$HOME/xbsapi.toml`
- `$PWD/xbsapi.toml`


### Database

xbsapi requires a database to store bookmarks. Supported database types are 
SQLite, PostgreSQL and MySQL. The database can be configured using the 
`XBSAPI_DATABASE_TYPE` and `XBSAPI_DATABASE_CONNECTION` env,
or the `Database.Type` and `Database.Connection` config properties.

**WARNING:** If you do not specify a database configuration, xbsapi will use
an in-memory SQLite database! As soon as xbsapi shuts down, all data
inside the in-memory database is gone!


#### SQLite File Example

```toml
[Database]
Type = "sqlite3"
Connection = "file:my-database.sqlite?cache=shared&_fk=1"
```


#### PostgreSQL Example *(using Docker for PostgreSQL)*

Run the database:

```sh
docker run -it --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=xbsapi \
  -p 127.0.0.1:5432:5432 \
  -d postgres:alpine
```

Configure `Database.Type` and `Database.Connection`:

```toml
[Database]
Type = "postgres"
Connection = "host=127.0.0.1 port=5432 dbname=xbsapi user=postgres password=postgres"
```


#### MySQL Example

```toml
[Database]
Type = "mysql"
Connection = "mysqluser:mysqlpassword@tcp(mysqlhost:port)/database?parseTime=true"
```


### Deployment

#### Custom

All that's needed is a [configuration](#configuration) and xbsapi can be
launched by e.g. running `./xbsapi` in a terminal.


#### Supervisor

To run xbsapi via `supervisord`, create a config like this inside
`/etc/supervisord.conf` or `/etc/supervisor/conf.d/xbsapi.conf`:

```ini
[program:xbsapi]
command=/path/to/binary/of/xbsapi
process_name=%(program_name)s
numprocs=1
directory=/home/xbsapi
autostart=true
autorestart=unexpected
startsecs=10
startretries=3
exitcodes=0
stopsignal=TERM
stopwaitsecs=10
user=xbsapi
redirect_stderr=false
stdout_logfile=/var/log/xbsapi.out.log
stdout_logfile_maxbytes=1MB
stdout_logfile_backups=10
stdout_capture_maxbytes=1MB
stdout_events_enabled=false
stderr_logfile=/var/log/xbsapi.err.log
stderr_logfile_maxbytes=1MB
stderr_logfile_backups=10
stderr_capture_maxbytes=1MB
stderr_events_enabled=false
```

**Note:** It is advisable to run xbsapi under its own, dedicated daemon
user (`xbsapi` in this example), so make sure to either adjust `directory`
as well as `user` or create a user called `xbsapi`.


#### OpenBSD rc

As before, create a configuration file under `/etc/xbsapi.toml`.

Then copy the [example rc.d script](examples/etc/rc.d/xbsapi) to
`/etc/rc.d/xbsapi` and copy the binary to e.g.
`/usr/local/bin/xbsapi`. Last but not least, update the `/etc/rc.conf.local`
file to contain the following line:

```conf
xbsapi_user="_xbsapi"
```

It is advisable to run xbsapi as a dedicated user, hence create the
`_xbsapi` daemon account or adjust the line above according to your setup.

You can now run xbsapi by enabling and starting the service:

```sh
rcctl enable xbsapi
rcctl start xbsapi
```


#### systemd

TODO


#### Docker

Official images are available on Docker Hub at 
[mrusme/xbsapi](https://hub.docker.com/r/mrusme/xbsapi) 
and can be pulled using the following command:

```sh
docker pull mrusme/xbsapi
```

GitHub release versions are available as Docker image tags (e.g. `1.0.0`). 
The `latest` image tag contains the latest code of the `master` branch.

It's possible to build xbsapi locally as a Docker container like this:

```sh
docker build -t xbsapi:latest . 
```

It can then be run using the following command:

```sh
docker run -it --rm --name xbsapi \
  -e XBSAPI_... \
  -e XBSAPI_... \
  -p 0.0.0.0:8000:8000 \
  xbsapi:latest
```

Alternatively a configuration TOML can be passed into the container like so:

```sh
docker run -it --rm --name xbsapi \
  -v /path/to/my/local/xbsapi.toml:/etc/xbsapi.toml \
  -p 0.0.0.0:8000:8000 \
  xbsapi:latest
```


#### Kubernetes

TODO


#### Render

Fork this repo into your GitHub account, adjust the
[`render.yaml`](render.yaml) accordingly and connect the forked repo [on
Render](https://dashboard.render.com/select-repo?type=blueprint).

Alternatively, you can also directly connect this public repo.


#### Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/mrusme/xbsapi)


#### DigitalOcean App Platform

[![Deploy to DO](https://www.deploytodo.com/do-btn-blue-ghost.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/mrusme/xbsapi/tree/master&refcode=9d48825ddae1)
  
Alternatively, fork this repo into your GitHub account, adjust the
[`.do/app.yaml`](.do/app.yaml) accordingly and connect the forked repo [on
DigitalOcean](https://cloud.digitalocean.com/apps/new).


#### DigitalOcean Function

Available soon.


#### Aamazon Web Services Lambda Function

TODO


#### Google Cloud Function

```sh
gcloud functions deploy GCFHandler --runtime go116 --trigger-http
```

TODO: Database


