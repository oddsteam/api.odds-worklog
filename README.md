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

## Host local
[http://localhost:8080/](http://localhost:8080/)

## API
local: http://localhost:8080/v1/
dev clound: http://worklog-dev.odds.team/api/v1/

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