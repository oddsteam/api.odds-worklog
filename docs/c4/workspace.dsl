workspace {
    model {
        odds_member = person "ODDS Member" {
            description "A member of ODDS-TEAM"
        }

        worklog = softwareSystem "ODDS-TEAM Worklog" {
            description "The worklog for ODDS members to record their man-day"

            web_app = container "Web Application" {
                description "The web application"
                technology "Angular"
            }

            api_app = container "API Application" {
                description "The API application"
                technology "Golang"
            }

            database = container "Database" {
                technology "MongoDB"
                tags "database"
            }
        }

        odds_member -> web_app "Uses"
        web_app -> api_app "Calls API from"
        api_app -> database "Read from and writes to"
    }

    views {
        systemContext worklog {
            include *
            autoLayout lr
        }

        container worklog {
            include *
            autoLayout lr
        }

        theme default

        styles {
            element "database" {
                shape Cylinder
            }
        }
    }
}
