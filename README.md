# API Odds-Worklog 

## Go version 1.11

## Run by docker-compose API+Mongodb

Run `docker-compose up --build -d`

If error is `The container name "/odds-worklog-mongo" is already.`

Run `docker stop odds-worklog-mongo` to stop container or

Run `docker rm odds-worklog-mongo` to remove old container

Then run `docker-compose up --build -d` again

If api not start because `Authentication failed` you can setup authen mongodb below.


## Setup authen mongodb

If first time, you must Run `docker volume create mongodbdata_odds_worklog` for create mongodb docker volume.

1. Start mongodb container <br> 
If container mongo name `odds-worklog-mongo` is not running <br>
Run `docker run -it --rm --name odds-worklog-mongo -d -p 27017:27017 -v mongodbdata_odds_worklog:/data/db mongo` <br>
If docker: Error response from daemon: Conflict. Run `docker rm odds-worklog-mongo` to remove old container <br>
If you have container using port `27017` you must be stop it and run command again.

2. Invoke mongodb container <br>
Run `docker exec -it odds-worklog-mongo bash`

3. Invoke mongodb <br>
Run `mongo`

4. Create user admin <br>
Run `use admin` <br>
If no user admin <br>
Run `db.createUser({user:"admin",pwd:"admin",roles:["root"]})` <br>
Else authen, run `db.auth('admin','admin')`

5. Create user for read/write data on `odds_worklog_db` <br>
Run `use odds_worklog_db` <br>
Run `db.createUser({user:"admin",pwd:"admin",roles:[{role:"readWrite",db:"odds_worklog_db"}]})`

6. Exit mongo and container<br>
Exit mongo in container, run `exit` <br>
Exit container, run `exit`


## Run by `go run main.go` <br>
Run `go run main.go` (at project path)


## API

local: http://localhost:8080/v1/

develop cloud: https://worklog-dev.odds.team/api/v1/

production cloud: https://worklog.odds.team/api/v1/

## Import mock data to mongodb

If you use to import data mock, data should be alive. <br>
Importion is optional.

* **Import user data** <br>
At project path<br>
```bash 
    mongoimport --host localhost --port 27017 --db odds_worklog_db --collection user --type json --file data/user.json --maintainInsertionOrder --jsonArray
```

* **Import site data** <br>
At project path<br>
```bash 
    mongoimport --host localhost --port 27017 --db odds_worklog_db --collection site --type json --file data/site.json --maintainInsertionOrder --jsonArray
```

* **Import customer data** <br>
At project path<br>
```bash 
    mongoimport --host localhost --port 27017 --db odds_worklog_db --collection customer --type json --file data/customer.json --maintainInsertionOrder --jsonArray
```

* **Import po data** <br>
At project path<br>
```bash 
    mongoimport --host localhost --port 27017 --db odds_worklog_db --collection po --type json --file data/po.json --maintainInsertionOrder --jsonArray
```

* **Import invoice data** <br>
At project path<br>
```bash 
    mongoimport --host localhost --port 27017 --db odds_worklog_db --collection invoice --type json --file data/invoice.json --maintainInsertionOrder --jsonArray
```

## Command go mockgen

GoMock is a mocking framework for the Go programming language.

[https://github.com/golang/mock](https://github.com/golang/mock)

[https://godoc.org/github.com/golang/mock/gomock](https://godoc.org/github.com/golang/mock/gomock)

user `mockgen -source="api/user/interface.go" -destination="api/user/mock/user_mock.go"`

income `mockgen -source="api/income/interface.go" -destination="api/income/mock/income_mock.go"`

login `mockgen -source="api/login/interface.go" -destination="api/login/mock/login_mock.go"`

file `mockgen -source="api/file/interface.go" -destination="api/file/mock/file_mock.go"`

site `mockgen -source="api/site/interface.go" -destination="api/site/mock/site_mock.go"`

customer `mockgen -source="api/customer/interface.go" -destination="api/customer/mock/customer_mock.go"`

po `mockgen -source="api/po/interface.go" -destination="api/po/mock/po_mock.go"`

invoice `mockgen -source="api/invoice/interface.go" -destination="api/invoice/mock/invoice_mock.go"`

### Swagger

After fill comments to each handler, you must be run `swag init` to generate swagger docs

[https://github.com/swaggo/swag](https://github.com/swaggo/swag)

[https://github.com/swaggo/echo-swagger](https://github.com/swaggo/echo-swagger)

local [http://localhost:8080/v1/swagger/index.html](http://localhost:8080/v1/swagger/index.html)

online [https://worklog-dev.odds.team/api/v1/swagger/index.html](http://worklog-dev.odds.team/api/v1/swagger/index.html)

### Run test

Run all test `go test ./...`

Run all test coverage by package `go test ./... -cover`

Run all test caoverage and view with html <br>
`go test -coverprofile=cover.out` <br>
`go test ./... -coverprofile=cover.out && go tool cover -html=cover.out`