## Simple bank (Work in progress)


### Working with database

### Tech stack

- golang
- postgres DB
- [golang/migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#linux-deb-package): data schema migration library
- [sqlc](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html#):An sql generation library

### Folder structure


#### commands

1. run postgres database

container name: postgres

port mapping: 5432

root user: root

```shell
docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres
```

2. Connect database in docker

```shell
docker exec -it postgres  psql -U root
```

3. Open terminal in docker

```shell
docker exec -it postgres /bin/sh
```

4. create a db after

```shell
createdb --username=root --owner=root simple_bank
```

5. database migration command

- path: path to migration file
- database: connected database
- verbose: up / down

```shell
migrate -path db/migration -database "postgresql://root:admin@localhost:5432/simple_bank?sslmode=disable" -verbose up
```
