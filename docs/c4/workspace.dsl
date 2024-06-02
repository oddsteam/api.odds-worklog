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

        google_signin = softwareSystem "Google Sign-In" {
            description "This helps us easily and securely sign in to our worklog app with our Google Account. It also manages the OAuth 2.0 flow and token lifecycle for us."
        }

        odds_member -> web_app "Uses"
        web_app -> api_app "Calls API from"
        api_app -> database "Read from and writes to"
        api_app -> google_signin "Integrates with"
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
