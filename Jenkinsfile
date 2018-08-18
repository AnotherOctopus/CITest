node {
        def app
        stage ('build') {
                try{    
                        checkout scm
                        sh 'echo ${PULLBRANCH}'
                        sh 'git fetch'
                        sh 'git checkout ${PULLBRANCH}' 
                        sh 'id'
                        app = docker.build("anotheroctopus/rovimage")
                }catch(error){
                        slackSend(color: "#00FF00",message: "Hum, we failed checking out the repo. Idk man")
                }
        }
        stage ('launchROV'){
                try{
                        sh 'docker login -u anotheroctopus -p 44Cobr@'
                        app.push(env.BRANCH_NAME)
                        sh 'id'
                        sh 'ssh -p 2112 sampi@dhtilly.ddns.net \'df -H\''
                        sh 'ssh -p 2112 sampi@dhtilly.ddns.net \'docker run anotheroctopus/rovimage:${PULLBRANCH}\''
                }catch(error){
                        slackSend(color: "#00FF00",message: "Launching the ROV failed. Probably some networking nonesense")
                }
        }
        stage('lint'){
                sh 'find . -iname "*.go" | xargs gofmt -d'
                sh 'find . -iname "*.py" | xargs pylint -d'
        }
        stage ('test'){
        }       
        stage ('post'){
                sh 'echo ${PULLMAKER}'
                sh 'echo ${REVIEWERS}'
                slackSend(color: "#00FF00",message: "Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Passed all Tests!")
                slackSend(color: "#00FF00",message: "Hey ${PULLMAKER}, you should bug ${REVIEWERS} to approve your pull")
        }
}
