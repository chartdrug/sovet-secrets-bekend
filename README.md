# Go RESTful API Starter Kit (Boilerplate)

[![GoDoc](https://godoc.org/github.com/qiangxue/sovet-secrets-bekend?status.png)](http://godoc.org/github.com/qiangxue/sovet-secrets-bekend)
[![Build Status](https://github.com/qiangxue/sovet-secrets-bekend/workflows/build/badge.svg)](https://github.com/qiangxue/sovet-secrets-bekend/actions?query=workflow%3Abuild)
[![Code Coverage](https://codecov.io/gh/qiangxue/sovet-secrets-bekend/branch/master/graph/badge.svg)](https://codecov.io/gh/qiangxue/sovet-secrets-bekend)
[![Go Report](https://goreportcard.com/badge/github.com/qiangxue/sovet-secrets-bekend)](https://goreportcard.com/report/github.com/qiangxue/sovet-secrets-bekend)

This starter kit is designed to get you up and running with a project structure optimized for developing
RESTful API services in Go. It promotes the best practices that follow the [SOLID principles](https://en.wikipedia.org/wiki/SOLID)
and [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). 
It encourages writing clean and idiomatic Go code. 

The kit provides the following features right out of the box:

* RESTful endpoints in the widely accepted format
* Standard CRUD operations of a database table
* JWT-based authentication
* Environment dependent application configuration management
* Structured logging with contextual information
* Error handling with proper error response generation
* Database migration
* Data validation
* Full test coverage
* Live reloading during development
 
The kit uses the following Go packages which can be easily replaced with your own favorite ones
since their usages are mostly localized and abstracted. 

* Routing: [ozzo-routing](https://github.com/go-ozzo/ozzo-routing)
* Database access: [ozzo-dbx](https://github.com/go-ozzo/ozzo-dbx)
* Database migration: [golang-migrate](https://github.com/golang-migrate/migrate)
* Data validation: [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
* Logging: [zap](https://github.com/uber-go/zap)
* JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)

## Getting Started

If this is your first time encountering Go, please follow [the instructions](https://golang.org/doc/install) to
install Go on your computer. The kit requires **Go 1.13 or above**.

[Docker](https://www.docker.com/get-started) is also needed if you want to try the kit without setting up your
own database server. The kit requires **Docker 17.05 or higher** for the multi-stage build support.

After installing Go and Docker, run the following commands to start experiencing this starter kit:

```shell
# download the starter kit
git clone https://github.com/qiangxue/sovet-secrets-bekend.git

cd sovet-secrets-bekend

# start a PostgreSQL database server in a Docker container
make db-start

# seed the database with some test data
make testdata

# run the RESTful API server
make run

# or run the API server with live reloading, which is useful during development
# requires fswatch (https://github.com/emcrisostomo/fswatch)
make run-live
```

At this time, you have a RESTful API server running at `http://127.0.0.1:8080`. It provides the following endpoints:

* `GET /healthcheck`: a healthcheck service provided for health checking purpose (needed when implementing a server cluster)
* `POST /v1/login`: authenticates a user and generates a JWT
* `GET /v1/albums`: returns a paginated list of the albums
* `GET /v1/albums/:id`: returns the detailed information of an album
* `POST /v1/albums`: creates a new album
* `PUT /v1/albums/:id`: updates an existing album
* `DELETE /v1/albums/:id`: deletes an album

Try the URL `http://localhost:8080/healthcheck` in a browser, and you should see something like `"OK v1.0.0"` displayed.

If you have `cURL` or some API client tools (e.g. [Postman](https://www.getpostman.com/)), you may try the following 
more complex scenarios:

```shell
# authenticate the user via: POST /v1/login
curl -X POST -H "Content-Type: application/json" -d '{"username": "demo", "password": "pass"}' http://localhost:8080/v1/login
# should return a JWT token like: {"token":"...JWT token here..."}

# with the above JWT token, access the album resources, such as: GET /v1/albums
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8080/v1/albums
# should return a list of album records in the JSON format
```

To use the starter kit as a starting point of a real project whose package name is `github.com/abc/xyz`, do a global 
replacement of the string `github.com/qiangxue/sovet-secrets-bekend` in all of project files with the string `github.com/abc/xyz`.


## Project Layout

The starter kit uses the following project layout:
 
```
.
├── cmd                  main applications of the project
│   └── server           the API server application
├── config               configuration files for different environments
├── internal             private application and library code
│   ├── album            album-related features
│   ├── auth             authentication feature
│   ├── config           configuration library
│   ├── entity           entity definitions and domain logic
│   ├── errors           error types and handling
│   ├── healthcheck      healthcheck feature
│   └── test             helpers for testing purpose
├── migrations           database migrations
├── pkg                  public library code
│   ├── accesslog        access log middleware
│   ├── graceful         graceful shutdown of HTTP server
│   ├── log              structured and context-aware logger
│   └── pagination       paginated list
└── testdata             test data scripts
```

The top level directories `cmd`, `internal`, `pkg` are commonly found in other popular Go projects, as explained in
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

Within `internal` and `pkg`, packages are structured by features in order to achieve the so-called
[screaming architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html). For example, 
the `album` directory contains the application logic related with the album feature. 

Within each feature package, code are organized in layers (API, service, repository), following the dependency guidelines
as described in the [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).


## Common Development Tasks

This section describes some common development tasks using this starter kit.

### Implementing a New Feature

Implementing a new feature typically involves the following steps:

1. Develop the service that implements the business logic supporting the feature. Please refer to `internal/album/service.go` as an example.
2. Develop the RESTful API exposing the service about the feature. Please refer to `internal/album/api.go` as an example.
3. Develop the repository that persists the data entities needed by the service. Please refer to `internal/album/repository.go` as an example.
4. Wire up the above components together by injecting their dependencies in the main function. Please refer to 
   the `album.RegisterHandlers()` call in `cmd/server/main.go`.

### Working with DB Transactions

It is the responsibility of the service layer to determine whether DB operations should be enclosed in a transaction.
The DB operations implemented by the repository layer should work both with and without a transaction.

You can use `dbcontext.DB.Transactional()` in a service method to enclose multiple repository method calls in
a transaction. For example,

```go
func serviceMethod(ctx context.Context, repo Repository, transactional dbcontext.TransactionFunc) error {
    return transactional(ctx, func(ctx context.Context) error {
        repo.method1(...)
        repo.method2(...)
        return nil
    })
}
```

If needed, you can also enclose method calls of different repositories in a single transaction. The return value
of the function in `transactional` above determines if the transaction should be committed or rolled back.

You can also use `dbcontext.DB.TransactionHandler()` as a middleware to enclose a whole API handler in a transaction.
This is especially useful if an API handler needs to put method calls of multiple services in a transaction.


### Updating Database Schema

The starter kit uses [database migration](https://en.wikipedia.org/wiki/Schema_migration) to manage the changes of the 
database schema over the whole project development phase. The following commands are commonly used with regard to database
schema changes:

```shell
# Execute new migrations made by you or other team members.
# Usually you should run this command each time after you pull new code from the code repo. 
make migrate

# Create a new database migration.
# In the generated `migrations/*.up.sql` file, write the SQL statements that implement the schema changes.
# In the `*.down.sql` file, write the SQL statements that revert the schema changes.
make migrate-new

# Revert the last database migration.
# This is often used when a migration has some issues and needs to be reverted.
make migrate-down

# Clean up the database and rerun the migrations from the very beginning.
# Note that this command will first erase all data and tables in the database, and then
# run all migrations. 
make migrate-reset
```

### Managing Configurations

The application configuration is represented in `internal/config/config.go`. When the application starts,
it loads the configuration from a configuration file as well as environment variables. The path to the configuration 
file is specified via the `-config` command line argument which defaults to `./config/local.yml`. Configurations
specified in environment variables should be named with the `APP_` prefix and in upper case. When a configuration
is specified in both a configuration file and an environment variable, the latter takes precedence. 

The `config` directory contains the configuration files named after different environments. For example,
`config/local.yml` corresponds to the local development environment and is used when running the application 
via `make run`.

Do not keep secrets in the configuration files. Provide them via environment variables instead. For example,
you should provide `Config.DSN` using the `APP_DSN` environment variable. Secrets can be populated from a secret
storage (e.g. HashiCorp Vault) into environment variables in a bootstrap script (e.g. `cmd/server/entryscript.sh`). 

## Deployment

The application can be run as a docker container. You can use `make build-docker` to build the application 
into a docker image. The docker container starts with the `cmd/server/entryscript.sh` script which reads 
the `APP_ENV` environment variable to determine which configuration file to use. For example,
if `APP_ENV` is `qa`, the application will be started with the `config/qa.yml` configuration file.

You can also run `make build` to build an executable binary named `server`. Then start the API server using the following
command,

```shell
./server -config=./config/prod.yml
```

```shell
chmod 400 sovet-ZEFNmBra.pem 
chmod 400 amaz.pem

-- vk old
docker save server:latest | bzip2 | ssh -i sovet-ZEFNmBra.pem centos@213.219.213.247 docker load
-- vk amster
docker save server:latest | bzip2 | ssh -i vk_amster.pem centos@89.208.219.91 docker load
-- amster
docker save server:latest | bzip2 | ssh -i amster.pem centos@89.208.219.143 docker load

-- amazon
docker save server:latest | bzip2 | ssh -i amaz.pem ec2-user@35.175.222.125 docker load

scp -i vk_amster.pem config/server-compose.yml centos@89.208.219.91:/home/centos 

-- vk
ssh -i sovet-ZEFNmBra.pem centos@213.219.213.247
-- vk amster
ssh -i vk_amster.pem centos@89.208.219.91
-- amster
ssh -i amster.pem ubuntu@89.208.219.143

-- amaz
ssh -i amaz.pem ec2-user@35.175.222.125

docker stop server 
docker rm server

docker run -it --rm -d -p 8080:8080 --name server server
docker run -it --rm -d -v /home/centos/logs:/var/log/app -p 8080:8080 --name server server 
docker stop server && docker rm server && docker run -it -d --restart unless-stopped -p 8080:8080 --name server server


 docker network create chart

ssh -i vk_amster.pem centos@89.208.219.91
docker stop centos_server_1 && docker rm centos_server_1 && docker-compose -f server-compose.yml up -d

rm /home/centos/logs/server.log
docker-compose -f server-compose.yml up -d

ssh -i vk_amster.pem centos@89.208.219.91 
docker stop centos_server_1 && docker rm centos_server_1 && docker-compose -f server-compose.yml up -d

docker-compose -f config/server-compose.yml up -d

docker-compose -f docker-compose.yml up


echo '' | sudo tee -a nginx.conf



scp -r -i vk_amster.pem config/nginx.conf centos@89.208.219.91:/home/centos/nginx 
scp -r -i vk_amster.pem config/privkey.pem centos@89.208.219.91:/home/centos/nginx 
scp -r -i vk_amster.pem config/fullchain.pem centos@89.208.219.91:/home/centos/nginx 

sudo cp /home/centos/nginx/* /etc/nginx

chmod g+x /home/centos && chmod g+x /home/centos/front

chmod +x /home/centos/front

sudo cp -r /home/centos/front/* /usr/share/nginx/html


cd /etc/nginx
sudo systemctl start nginx
sudo systemctl stop nginx
sudo systemctl reload nginx
sudo systemctl status nginx

sudo systemctl stop httpd.service
sudo systemctl status httpd.service
sudo systemctl start httpd.service
sudo systemctl reload httpd.service

cd /etc/httpd/conf.d
sudo nano yourDomainName.conf

https://www.digitalocean.com/community/tutorials/how-to-secure-apache-with-let-s-encrypt-on-centos-7

перевыпуск серта 
sudo certbot --nginx -d chartdrug.com -d www.chartdrug.com
udo cp /etc/letsencrypt/live/chartdrug.com/fullchain.pem /etc/nginx/fullchain.pem
sudo cp /etc/letsencrypt/live/chartdrug.com/privkey.pem /etc/nginx/privkey.pem
sudo systemctl stop nginx
sudo systemctl start nginx




sudo systemctl stop nginx
sudo systemctl start httpd.service
sudo certbot --apache -d chartdrug.com -d www.chartdrug.com -d www.calcpharm.com -d calcpharm.com
//sudo certbot certonly --manual -d chartdrug.com -d www.chartdrug.com -d www.calcpharm.com -d calcpharm.com
certbot --nginx -d chartdrug.com -d www.chartdrug.com
sudo systemctl stop httpd.service && sudo systemctl start nginx

серты тут 12/09/2022
/etc/letsencrypt/live/chartdrug.com/fullchain.pem
Your key file has been saved at:
/etc/letsencrypt/live/chartdrug.com/privkey.pem

// logs
docker exec -it server bash
docker exec -it centos_server_1 bash

tail -n 10 /home/centos/logs/server.log

sudo cd /var/log/nginx

docker images --all
docker logs server


// POSTGRES_DB
docker logs habr-pg-13.3
docker exec -it habr-pg-13.3 bash

установка постгри в контейнере 
https://habr.com/ru/post/578744/

//docker run --name habr-pg-13.3 -p 5432:5432 -e POSTGRES_USER=habrpguser -e POSTGRES_PASSWORD=pgpwd4habr -e POSTGRES_DB=habrdb -e PGDATA=/var/lib/postgresql/data/pgdata -d -v "/absolute/path/to/directory-with-data":/var/lib/postgresql/data -v "/absolute/path/to/directory-with-init-scripts":/docker-entrypoint-initdb.d postgres:13.3
docker run --name habr-pg-13.3 --restart unless-stopped -p 5432:5432 -e POSTGRES_USER=chatrdruguser -e POSTGRES_PASSWORD=pgpwd4chatrdrug -e POSTGRES_DB=chatrdrug -e PGDATA=/var/lib/postgresql/data/pgdata -d -v "/absolute/path/to/directory-with-data":/var/lib/postgresql/data postgres:13.3
docker run --name habr-pg-13.3-2 --restart unless-stopped -p 5432:5432 -e POSTGRES_USER=chatrdruguser -e POSTGRES_PASSWORD=pgpwd4chatrdrug -e POSTGRES_DB=chatrdrug -e PGDATA=/var/lib/postgresql/data/pgdata -d -v "/home/centos/pg":/var/lib/postgresql/data postgres:13.3
docker run -it -d --restart unless-stopped -p 5432:5432 --name habr-pg-13.3 postgres:13.3

//docker run --name habr-pg-13.3 -p 5432:5432 -e POSTGRES_USER=habrpguser -e POSTGRES_PASSWORD=pgpwd4habr -e POSTGRES_DB=habrdb -e PGDATA=/var/lib/postgresql/data/pgdata -d -v "$(pwd)":/var/lib/postgresql/data -v "$(pwd)/../2. Init Database":/docker-entrypoint-initdb.d postgres:13.3

docker run --restart always --name mail -p 587:587 -e RELAY_HOST=smtp.chartdrug.com -e RELAY_PORT=587 -e RELAY_USERNAME=info@chartdrug.com -e RELAY_PASSWORD=secretpassword -d bytemark/smtp

docker-compose -f pg-compose.yml up


scp -i vk_amster.pem config/pg-compose.yml centos@89.208.219.91:/home/centos/pg 


данны бд
/absolute/path/to/directory-with-data

[centos@sovet directory-with-data]$ sudo cp -R pgdata /home/centos/pg/pgdata


/home/centos/pg


утсановка nginx

sudo dnf install nginx
sudo systemctl enable nginx
sudo systemctl start nginx



sudo yum install -y yum-utilssudo dnf install yum

// kafka
https://developer.confluent.io/quickstart/kafka-docker/

scp -i vk_amster.pem config/kafka-compose.yml centos@89.208.219.91:/home/centos/kafka 
docker stop zookeeper && docker rm zookeeper && docker stop broker && docker rm broker 
docker-compose -f kafka/kafka-compose.yml up -d && docker restart centos_server_1
новая кафка без зукипера
scp -i vk_amster.pem config/kafka2-compose.yml centos@89.208.219.91:/home/centos/kafka 

docker-compose -f kafka-compose.yml up -d

новая кафка без зукипера
docker-compose -f kafka/kafka2-compose.yml up -d

docker exec broker \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic calc_injection
             
docker exec kafka_kafka1_1 \
kafka-topics --bootstrap-server broker:9092 \
             --create \
             --topic calc_injection
             
// если докер не стартует             
sudo dockerd

установить curl
apk add curl

список текущих портов
sudo ss -tlpuna


```

стопнуть процесс
sudo lsof -nP -i4TCP:8080 | grep LISTEN

main    18975 kdereshev   12u  IPv6 0x4adf38b374f1995      0t0  TCP *:8080 (LISTEN) 

sudo kill -9 18975