workspace {
    model {
        odds_member = person "ODDS Member" {
            description "A member of ODDS-TEAM"
        }

        odds_admin = person "ODDS Admin" {
            description "A admin of ODDS Worklog"
        }

        worklog = softwareSystem "ODDS-TEAM Worklog" {
            description "The worklog for ODDS members to record their man-day"

            ios_app = container "iOS App" {
                tags "mobile"
            }

            android_app = container "Android App" {
                tags "mobile"
            }

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
        odds_member -> ios_app "Uses"
        odds_member -> android_app "Uses"
        odds_admin -> web_app "Exports income and deletes inactive members from"
        web_app -> api_app "Calls API from"
        ios_app -> api_app "Calls API from"
        android_app -> api_app "Calls API from"
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

            element "mobile" {
                shape MobileDevicePortrait
            }
        }
    }
}
