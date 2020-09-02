
pipeline {
  agent any
  stages {
    stage('registry') {
      steps {
        script {
            sh "echo $GIT_BRANCH" 
            def shortenv = "dev"
            if (GIT_BRANCH == "origin/new-develop") {
                shortenv = 'dev'
            } else if (GIT_BRANCH == "origin/master"){
                shortenv = 'prod'
            }
            withDockerRegistry(credentialsId: 'ecr:ap-southeast-1:registry-aws', url: 'https://844772501268.dkr.ecr.ap-southeast-1.amazonaws.com') {
                        def worklogAPI = docker.build("844772501268.dkr.ecr.ap-southeast-1.amazonaws.com/api-worklog:$shortenv")
                        
                        //push image
                        worklogAPI.push()
            }
        }
      }
    }
  }
}