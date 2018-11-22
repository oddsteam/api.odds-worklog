# API Odds-Worklog 

## Go version 1.11

## Run Docker API+MONGODB<br>
Run `docker-compose up --build -d`<br>
*Note:* If you add new 3rd party package, you must run `dep ensure` for setup dependency.

## Setup run by `go run main.go` <br>
* **Docker mongodb**<br>
Run `docker run -it -d -p 27017:27017 mongo`

* **Setup Authen mongodb**<br>
1. In .env file, change `MONGO_DB_HOST = "mongodb:27017"` `mongodb` to `localhost` <br>
*Note:* Don't `commit` this file (.env)

2. Run `docker exec -it CONTAINER_MONGODB_NAME bash`
<br>get `CONTAINER_MONGODB_NAME` from `docker ps` NAMES

3. Run `mongo`

4. Run `use odds_worklog_db`

5. Run `db.createUser({user:"admin",pwd:"admin",roles:[{role:"readWrite",db:"odds_worklog_db"}]})`

* **Setup dependency**<br>
dep is a dependency management tool for Go <br>
Run `go get -u github.com/golang/dep/cmd/dep` <br>
Run `dep ensure` (at project path)

* **Import user data** <br>
At project path<br>
```bash 
    mongoimport --host localhost --port 27017 --db odds_worklog_db --collection user --type json --file user.json --maintainInsertionOrder --jsonArray
```

* **Run `go run main.go` (at project path)**


## Set up Swagger
[https://github.com/swaggo/echo-swagger](https://github.com/swaggo/echo-swagger)

After fill Comment to each handler, you must run `swag init` to generate docs swagger 

## Host local
[http://localhost:8080/](http://localhost:8080/)

## Command go mockgen
GoMock is a mocking framework for the Go programming language.

[https://github.com/golang/mock](https://github.com/golang/mock)

[https://godoc.org/github.com/golang/mock/gomock](https://godoc.org/github.com/golang/mock/gomock)

user `mockgen -source="api/user/interface.go" -destination="api/user/mock/user_mock.go"`

income `mockgen -source="api/income/interface.go" -destination="api/income/mock/income_mock.go"`

login `mockgen -source="api/login/interface.go" -destination="api/login/mock/login_mock.go"`

## API
local: http://localhost:8080/v1/

dev clound: http://worklog-dev.odds.team/api/v1/

### Swagger
local [http://localhost:8080/v1/swagger/index.html](http://localhost:8080/v1/swagger/index.html)

online [http://worklog-dev.odds.team/api/v1/swagger/index.html](http://worklog-dev.odds.team/api/v1/swagger/index.html)

### User
| Method    | Path          |
| ---       | ---           |
| GET       | /users        |
| GET       | /users/:id    |
| POST      | /users/:id    |
| POST      | /login        |
| PUT       | /users/:id    |
| PATCH     | /users/:id    |
| DELETE    | /users/:id    |

### Income
| Method    | Path              |
| ---       | ---               |
| GET       | /incomes/status   |
| POST      | /incomes/         |
| PUT       | /incomes/:id      |   