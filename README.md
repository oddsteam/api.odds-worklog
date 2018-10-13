# API Odds-Worklog 

## Setup
* **Install dep**, dep is a dependency management tool for Go <br>
Run `go get -u github.com/golang/dep/cmd/dep`

* **Setup dependency**  <br>
Run `dep ensure`

* **Run Docker** <br>
Run `docker-compose up --build -d`

* **Import user data** <br>
```bash 
    mongoimport --host localhost --port 27017 --db odds_worklog_db --collection user --type json --file user.json --maintainInsertionOrder --jsonArray
```

## Host local
[http://localhost:8080/](http://localhost:8080/)

## API
GET /user

GET /user/:id

POST /user

PUT /user

DELETE /user/:id