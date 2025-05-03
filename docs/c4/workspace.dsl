workspace {
    model {
        odds_member = person "ODDS Member" {
            description "A member of ODDS-TEAM"
        }

        odds_admin = person "ODDS Admin" {
            description "A admin of ODDS Worklog"
        }

        central_queue = softwareSystem "Central Queue" {
            description "Central RabbitMQ works as Enterprise System Bus"
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

            friends_log_worker = container "FriendsLog Worker" {
                description "The Subscriber of FriendsLog events"
                technology "Golang"
            }

            database = container "Database" {
                technology "MongoDB"
                tags "database"
            }

            group "Network" {
                elb = container "ELB" {
                    description "shared ELB for HTTPS"
                }

                reverse_proxy = container "Reverse Proxy" {
                    description "Separating prod and dev environment requests"
                }

                local_nginx = container "Local Reverse Proxy" {
                    description "Separating web and api requests"
                }
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
        friends_log_worker -> central_queue "listen incomes_created and incomes_updated events"
        friends_log_worker -> database "update incomes from FriendsLog"



        odds_member -> elb "Sends HTTPS requests"
        elb -> reverse_proxy "HTTP"
        reverse_proxy -> local_nginx
        local_nginx -> web_app
        local_nginx -> api_app

        deployment = deploymentEnvironment "worklog-dev" {
            deploymentNode "Huawei Cloud" {
                containerInstance elb
                deploymentNode "Worklog instance" {
                    containerInstance reverse_proxy
                    deploymentNode "worklog-web-dev" {
                        containerInstance local_nginx
                        containerInstance web_app
                        containerInstance api_app
                    }
                }
            }
        }
    }

    views {
        systemContext worklog {
            include *
            autoLayout lr
        }

        container worklog {
            include odds_member
            include odds_admin
            include ios_app
            include android_app
            include web_app
            include api_app
            include friends_log_worker
            include database
            autoLayout lr
        }

        deployment * deployment {
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
