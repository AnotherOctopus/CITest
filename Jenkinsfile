node {
        def app
        stage ('build') {
                try{    
                        checkout scm
                        sh 'echo ${PULLBRANCH}'
                        sh 'git fetch'
                        sh 'git checkout ${PULLBRANCH}' 
                }catch(error){
                        slackSend(color: "#FF0000",message: "Hum, we failed checking out the repo. Idk man")
                        error("SOURCE FAILED")
                }
                try{
                        app = docker.build("anotheroctopus/rovimage")
                }catch(error){
                        slackSend(color: "#FF0000",message: "So the docker image didn't build, so its either Scotty's fault or the Dockerfile")
                        error("BUILD FAILED")
                }
        }
        stage ('launchROV'){
                try{
                        sh 'docker login -u anotheroctopus -p 44Cobr@'
                        app.push(env.BRANCH_NAME)
                        sh 'id'
                        //sh 'ssh -p 2112 sampi@dhtilly.ddns.net \'df -H\''
                        //sh 'ssh -p 2112 sampi@dhtilly.ddns.net \'docker run anotheroctopus/rovimage:${PULLBRANCH}\''
                }catch(error){
                        slackSend(color: "#FF0000",message: "Launching the ROV failed. Probably some networking nonesense")
                        error("ROV FAILED TO LAUNCH")
                }
        }
        stage('lint'){
                try{
                        //golint = sh(returnStdout:true, script: 'find . -iname "*.go" | xargs gofmt -d').trim()
                }catch(error){
                        slackSend(color: "#FF0000",message: "Linting Go Files Failed!")
                        error("LINT FAILED")
                }
                try{
                        //pylint = sh(returnStdout:true, script: 'find . -iname "*.py" | xargs pylint -d').trim()
                }catch(error){
                        slackSend(color: "#FF0000",message: "Linting Python Files Failed!")
                        error("LINT FAILED")
                }
                try{
                        //eslint = sh(returnStdout:true, script: 'find . -iname "*.jsx" | xargs eslint -d').trim()
                }catch(error){
                        slackSend(color: "#FF0000",message: "Linting React Files Failed!")
                        error("LINT FAILED")
                } 
        }
        stage ('test'){
                sh 'ls .'
                sh 'python tests/testem.py'
        }       
        stage ('post'){
                slackSend(color: "#00FF00",message: "Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Passed all Tests!")
                slackSend(color: "#00FF00",message: "Hey ${PULLMAKER}, you should bug ${REVIEWERS} to approve your pull")
        }
}
