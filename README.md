# API Odds-Worklog 

## Go version 1.11

## First step, please setup authen mongodb

1. Run `docker run -it --name mongodb -d -p 27017:27017 mongo` <br> 
If docker: Error response from daemon: Conflict. Run `docker rm $(docker ps -a -q)` for remove history containers are run. Or you can rename this container to other name.

2. Run `docker exec -it mongodb bash`

3. Run `mongo` or `mongo -u admin -p admin --authenticationDatabase admin`

4. Run `use odds_worklog_db`

5. Run `db.createUser({user:"admin",pwd:"admin",roles:[{role:"readWrite",db:"odds_worklog_db"}]})`

## Run by docker-compose API+Mongodb

Run `docker-compose up --build -d`<br>
*Note:* If you add new 3rd party package, you must be run `dep ensure` for setup dependency.

## Run by `go run main.go` <br>
dep is a dependency management tool for Go

1. Run `go get -u github.com/golang/dep/cmd/dep`

2. Run `dep ensure` (at project path)

3. In .env file, change `MONGO_DB_HOST = "mongodb:27017"` from `mongodb` to `localhost` <br>
*Note:* Don't `commit` this file (.env)

4. Run `go run main.go` (at project path)

## API

local: http://localhost:8080/v1/

develop cloud: https://worklog-dev.odds.team/api/v1/

production cloud: https://worklog.odds.team/api/v1/

## Import mock data to mongodb
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

po `mockgen -source="api/po/interface.go" -destination="api/po/mock/po_mock.go"`

invoice `mockgen -source="api/invoice/interface.go" -destination="api/invoice/mock/invoice_mock.go"`

### Swagger

After fill comments to each handler, you must be run `swag init` to generate swagger docs

[https://github.com/swaggo/swag](https://github.com/swaggo/swag)

[https://github.com/swaggo/echo-swagger](https://github.com/swaggo/echo-swagger)

local [http://localhost:8080/v1/swagger/index.html](http://localhost:8080/v1/swagger/index.html)

online [https://worklog-dev.odds.team/api/v1/swagger/index.html](http://worklog-dev.odds.team/api/v1/swagger/index.html)