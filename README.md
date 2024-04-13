# ODDS Worklog API

ODDS Worklog API written in Go

**Table of contents**

- [Requirements](#requirements)
- [Getting started](#getting-started)
- [Contributing](#contributing)

## Requirements

ODDS Worklog API is tested with:

|         | Main version |
| ------- | ------------ |
| Go      | 1.13         |
| MongoDB | 4.0.3        |

## Getting started

### Running in Docker

Install [Docker Community Edition (CE)](https://docs.docker.com/engine/install/) on your machine.

```sh
docker-compose up --build -d
```

If you found the error "The container name "/odds-worklog-mongo" is already.", please follow the steps below.

1. Run `docker stop odds-worklog-mongo` to stop container or `docker rm odds-worklog-mongo` to remove old container
1. Run `docker-compose up --build -d` again

If the API does not start due to the error "Authentication failed", please see the section below.

### Setup authen mongodb

If first time, you must Run `docker volume create mongodbdata_odds_worklog` for create mongodb docker volume.

1. Start mongodb container <br>
   If container mongo name `odds-worklog-mongo` is not running <br>
   Run `docker-compose -f docker-compose.local.yaml up -d mongodb` <br>
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

### Run by `go run main.go` <br>

Run `go run main.go` (at project path)

### API

local: http://localhost:8080/v1/

develop cloud: https://worklog-dev.odds.team/api/v1/

production cloud: https://worklog.odds.team/api/v1/

### Import mock data to mongodb

If you use to import data mock, data should be alive. <br>
Importion is optional.

At project path<br>

```bash
    ./scripts/import_all_data
```

### Command go mockgen

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

Run all test `./runtests/README.md`

Run all test coverage by package `go test ./... -cover`

Run all test coverage and view with html <br>
`go test -coverprofile=cover.out` <br>
`go test ./... -coverprofile=cover.out && go tool cover -html=cover.out`

## Cannot login worklog-dev

Use the command below to reset dev database.

`ssh worklog-huawei docker exec -t mongodb-dev  mongorestore -u admin -p admin -d odds_worklog_db  ./data/odds_worklog_db`

## Running Get Student Loan Script

`scripts/get_student_loan.go` requires JSESSIONID and X-CSRF-TOKEN to run. Get it by logging in at the student loan website. (You will know how to do that if you have access to it.)

#### On Local

```
SESSION="JSESSIONID=DXpquTUIAivMunzsgHls28n3FBxh-p7ECDeGLijW.node1" CSRF="3cc94f0d-e690-4fe1-89f8-1e6c2a51d5bb" go run scripts/get_student_loan.go
```

#### On Dev

```
ssh worklog-huawei docker exec -e SESSION="JSESSIONID=dwEi2vEj0qZM5-KNIK5xWfzFABFVKUBF7K8m8T37.node1" -e CSRF="1692a557-2606-4e7e-8911-bf0b0b5e481d" worklog-api-dev ./get_student_loan
```

## Contributing

Want to help build ODDS Worklog API?

1. Fork it (https://github.com/oddsteam/api.odds-worklog/fork)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request
