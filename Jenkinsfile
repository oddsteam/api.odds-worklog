
pipeline {
  agent any
  stages {
    stage('registry') {
      steps {
        script {
            sh "echo $env.BRANCH_NAME" 
            def shortenv = "dev"
            if (env.BRANCH_NAME == "develop") {
                shortenv = "dev"
            } else if (env.BRANCH_NAME == "master") {
                shortenv = "prod"
            }
            withDockerRegistry(credentialsId: 'ecr:ap-southeast-1:sapanboon-aws', url: 'https://844772501268.dkr.ecr.ap-southeast-1.amazonaws.com') {
                        docker.image("api-worklog:dev").push()
            }
        }
      }
    }
  }
}