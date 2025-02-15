## Simple bank (WIP)


### Tech stack

#### Design tool
- [dbdiagram](https://dbdiagram.io/home): It is used to design data table schema and support to export sql.

#### Dev Stack
- golang
- postgres DB
- [golang/migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#linux-deb-package): data schema migration library
- [sqlc](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html#):An sql generation library
- [viper](https://github.com/spf13/viper): environment configuration
- [goMock](https://github.com/golang/mock): mock database, used for test

#### DB migration
Create a new migration

1. Create a new migrtion and add sql in generated migration file
```shell
migrate create -ext sql -dir db/migration -seq <migration_name>
```
2. Run migration

```shell
make migrateup
```
3. Create query file in `db/query` folder

4. Generate query file in go

```shell
make sqlc
```

5. Re-generate mock file for test

```shell
make mock
```

6. Run test to make sure all existed test pass

```shell
make test
```

#### Test library
- [Testify](https://github.com/stretchr/testify): used to assert test result?

#### Mock install
After install, it needs to setting path. 
```shell
vi ~/.zshrc

# In the file, write below line, then save and exit
# export PATH=$PATH:~/go/bin

source ~/.zshrc

```

```shell
go install github.com/golang/mock/mockgen@v1.6.0
```
Generate mocked file/interface, first parameter is directory, second parameter is interface.

```shell
mockgen -destination db/mock/store.go github.com/liang3030/simple-bank/db/sqlc IStore
```
  

### Folder structure
```text

simple-bank/
├── db
│   ├── migration
│   │   ├── **_schema.down.sql
│   │   ├── **_schema.up.sql
│   ├── query
│   │   ├── **.sql
│   ├── sqlc
│   │   ├── **.go
├── util
│   ├── random.go
├── main.go
├── exported.sql 
├── go.mod
├── go.sum
├── Makefile
├── sqlc.yaml
├── README.md
└── .gitignore

```text

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

### Deployment

build docker image
```shell
docker build -t simple-bank:latest .
```

start container
```shell
docker run --name simpleBank -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:admin@XXX.XXX.XXX:5432/simple_bank?sslmode=disable" simple-bank:latest
```

create a network
```shell
docker network create bank-network
```

connect to a network 
```shell
docker network connect bank-network ${container name: postgres}
```

start container by conneted network 
```shell
docker run --name simpleBank --netowrk bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:admin@postgres:5432/simple_bank?sslmode=disable" simple-bank:latest
```

### Others

#### openssl
generate random string of 64 characters, and only takes 32 characters
```shell
 openssl rand -hex 64 | head -c 32
```
#### aws command
get secrect from aws
```shell
 aws secretsmanager get-secret-value --secret-id ${secret name: simple_bank} --query ${extract key: SecretString}
```

get secrect key and value, write it to target file

```shell
aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]'> app.env
```

login to aws ecr

```shell
aws ecr get-login-password | docker login --username AWS --password-stdin 679755225703.dkr.ecr.eu-central-1.amazonaws.com
```
update cubectl configuration from aws
```shell
 aws eks update-kubeconfig --name simple-bank --region eu-central-1
```

#### kubectl
apply aws auth file
```shell
kubectl apply -f eks/aws-auth.yml
```

tool
(k9s)[https://k9scli.io/]

#### jq
[jq](https://jqlang.github.io/jq/ ) is a lightweight and flexible command-line JSON processor.

#### cert manager

#### SSL/TLS


### Links
[Ingress nginx](https://github.com/kubernetes/ingress-nginx)

Error message: - Find out reason

```shell
interface conversion: interface {} is nil, not string
/usr/local/go/src/runtime/iface.go:275 (0x31e7ba4)
	panicdottypeE: panic(&TypeAssertionError{iface, have, want, ""})
/Users/liangzhang/Documents/programming_lang/simple-bank/db/sqlc/store.go:68 (0x35ee48f)
	(*SQLStore).TransferTx.func1: txName := ctx.Value(txKey).(string)
/Users/liangzhang/Documents/programming_lang/simple-bank/db/sqlc/store.go:35 (0x35ed766)
	(*SQLStore).execTx: err = fn(q)
/Users/liangzhang/Documents/programming_lang/simple-bank/db/sqlc/store.go:66 (0x35ed92a)
	(*SQLStore).TransferTx: err := store.execTx(ctx, func(q *Queries) error {
/Users/liangzhang/Documents/programming_lang/simple-bank/api/transfer.go:37 (0x35f0b79)
	(*Server).transfer: result, err := server.store.TransferTx(ctx, arg)
/Users/liangzhang/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/context.go:185 (0x35e68ae)
	(*Context).Next: c.handlers[c.index](c)
/Users/liangzhang/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/recovery.go:102 (0x35e689b)
	CustomRecoveryWithWriter.func1: c.Next()
/Users/liangzhang/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/context.go:185 (0x35e59e4)
	(*Context).Next: c.handlers[c.index](c)
/Users/liangzhang/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/logger.go:249 (0x35e59cb)
	LoggerWithConfig.func1: c.Next()
/Users/liangzhang/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/context.go:185 (0x35e4dd1)
	(*Context).Next: c.handlers[c.index](c)
/Users/liangzhang/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/gin.go:633 (0x35e4840)
	(*Engine).handleHTTPRequest: c.Next()
/Users/liangzhang/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/gin.go:589 (0x35e44d1)
	(*Engine).ServeHTTP: engine.handleHTTPRequest(c)
/usr/local/go/src/net/http/server.go:3210 (0x344fced)
	serverHandler.ServeHTTP: handler.ServeHTTP(rw, req)
/usr/local/go/src/net/http/server.go:2092 (0x34462ef)
	(*conn).serve: serverHandler{c.server}.ServeHTTP(w, w.req)
/usr/local/go/src/runtime/asm_amd64.s:1700 (0x324fa40)
	goexit: BYTE	$0x90	// NOP
```