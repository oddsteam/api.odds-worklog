# API Odds-Worklog 

## Setup
* **Install dep**, dep is a dependency management tool for Go <br>
Run `go get -u github.com/golang/dep/cmd/dep`

* **Setup dependency**  <br>
Run `dep ensure`

* **Run Docker** <br>
Run `docker-compose up --build -d` <br>
*Note:* If you add new 3rd party package, you must run `dep ensure` for setup dependency.

* **Import user data** <br>
```bash 
    mongoimport --host localhost --port 27017 --db odds_worklog_db --collection user --type json --file user.json --maintainInsertionOrder --jsonArray
```

## Set up Swagger
[https://github.com/swaggo/echo-swagger](https://github.com/swaggo/echo-swagger)

After fill Comment to each handler, you must run `swag init` to generate docs swagger 

## Host local
[http://localhost:8080/](http://localhost:8080/)

## Command go mockgen
GoMock is a mocking framework for the Go programming language.

[https://github.com/golang/mock](https://github.com/golang/mock)

[https://godoc.org/github.com/golang/mock/gomock](https://godoc.org/github.com/golang/mock/gomock)

user `mockgen -source="api/user/interface.go" -destination="api/user/mock/income_mock.go"`

income `mockgen -source="api/income/interface.go" -destination="api/income/mock/income_mock.go"`

login `mockgen -source="api/login/interface.go" -destination="api/login/mock/income_mock.go"`

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