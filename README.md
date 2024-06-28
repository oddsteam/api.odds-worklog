# ODDS Worklog API

ODDS Worklog API written in Go

**Table of Contents**

- [Requirements](#requirements)
- [Getting started](#getting-started)
- [Running Worklog API in Docker](#running-worklog-api-in-docker)
- [Contributing](#contributing)

## Requirements

ODDS Worklog API is tested with:

|         | Main version |
| ------- | ------------ |
| Go      | 1.13         |
| MongoDB | 4.0.3        |

## Architecture

See our [C4 Model](./docs/c4/).

## Getting Started

### Running Worklog API in Docker

1. Install [Docker Community Edition (CE)](https://docs.docker.com/engine/install/) on your machine.
1. Run the following command:

    ```sh
    docker compose up --build -d
    ```

If the API does not start due to the error "Authentication failed", please see the [MongoDB Authentication Setup](#mongodb-authentication-setup) section below.

#### MongoDB Authentication Setup

1. If the MongoDB container hasn't been started yet, run the following command:

   ```bash
   docker compose -f docker-compose.local.yaml up -d mongodb
   ```

1. Enter the MongoDB container's shell using the command below:

   ```bash
   docker compose -f docker-compose.local.yaml exec mongodb bash
   ```

1. After that, we'll invoke MongoDB shell by

   ```bash
   mongo
   ```

1. To authenticate an admin user, run the commands below:

   ```bash
   use admin
   db.auth("admin", "admin")
   ```

   If you found an error "Error: Authentication failed", it means that the admin user probably doesn't exist yet, we need to create it first by running the following command:

   ```bash
   db.createUser({user: "admin", pwd: "admin", roles:["root"]})
   ```

1. Create an user to manage data on the `odds_worklog_db` database, run the following commands:

   ```bash
   use odds_worklog_db
   db.createUser({user: "admin", pwd: "admin", roles: [{role: "readWrite",db: "odds_worklog_db"}]})
   ```

1. Now we can exit the MongoDB container:

   ```bash
   exit
   ```

### Starting Worklog API on Local Machine

Run the following command at the project path.

```bash
go run main.go
````

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

Run all test `./runtests.md`

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
