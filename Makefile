createdb:
#sqlserver
	sqlcmd -S localhost -U sa -P 1230 -Q "CREATE DATABASE simple_bank;"
dropdb:
#sqlserver
	sqlcmd -S localhost -U sa -P 1230 -Q "DROP DATABASE simple_bank;"
add-migrate:
#	migrate -database "sqlserver://sa:1230@localhost:1433?database=simple_bank&sslmode=disable" -path db/migration -verbose up
	migrate -database "mysql://root:1230@tcp(localhost:3306)/simple_bank?charset=utf8&parseTime=True&loc=Local" -path db/migration -verbose up 1

drop-migrate:
#	migrate -database "sqlserver://sa:1230@localhost:1433?database=simple_bank&sslmode=disable" -path db/migration -verbose down
	migrate -database "mysql://root:1230@tcp(localhost:3306)/simple_bank?charset=utf8&parseTime=True&loc=Local" -path db/migration -verbose down 1
force:
#	migrate -database "sqlserver://sa:1230@localhost:1433?database=simple_bank&sslmode=disable" -path db/migration force 1
	migrate -database "mysql://root:1230@tcp(localhost:3306)/simple_bank?charset=utf8&parseTime=True&loc=Local" -path db/migration force $(v)
create-migrate:
	migrate create -ext sql -dir db/migration -seq init_schema

#add-migration migrationname: 根据传入的参数创建新的迁移文件
#drop-migration migrationname: 根据传入的参数删除迁移文件
create-migration:
	migrate create -ext sql -dir db/migration -seq $(name)
delete-migration:
	rm db/migration/$(name)


#docker-mysql
createmysql:
	docker run --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=1230 -d mysql:latest
execmysql:
	docker exec -it mysql bash
logsmysql:
	docker logs mysql

#sqlc
sqlc:
	sqlc generate

#test
test:
	go test -v -cover ./...

#mock
mock:
	mockgen -destination=db/mock/store.go -package=mockdb github.com/ShenPingYuan/go-webdemo/db/sqlc Store

server:
	go run main.go

#docker
docker-build:
	docker build -t simplebank:latest .

.PHONY: createdb dropdb add-migrate drop-migrate force create-migrate createmysql execmysql logsmysql sqlc test server

