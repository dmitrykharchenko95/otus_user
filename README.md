# otus_user

## config
```shell
export SVC_ADDRESS=:8000
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_USER=postgres
export DB_PASS=postgres
export DB_NAME=main_db
```

## build docker image
```shell
GOOS=linux GOARCH=amd64 go build -o ./docker/bin/otus_user cmd/app/main.go 
docker build -t dmitrykharchenko95/otus_user:v0.0.8  ./docker

```