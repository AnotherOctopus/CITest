def WindDown(errorname){
        sh 'deactivate'
        error(errorname)
}
node {
        def app
        stage ('setup_virtualenv'){
                sh 'pwd'
                sh 'virtualenv venv'
                sh '. venv/bin/activate'
                sh 'which python'
                sh 'pip install pylint'
        }
        stage ('build') {
                try{    
                        checkout scm
                        sh 'echo ${PULLBRANCH}'
                        sh 'git checkout ${PULLBRANCH}' 
                        sh 'git pull'
                }catch(error){
                        slackSend(color: "#FF0000",message: "Hum, we failed checking out the repo. Idk man")
                        WindDown("SOURCE FAILED")
                }
                try{
                        app = docker.build("anotheroctopus/rovimage")
                }catch(error){
                        slackSend(color: "#FF0000",message: "So the docker image didn't build, so its either Scotty's fault or the Dockerfile")
                        WindDown("BUILD FAILED")
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
                        WindDown("ROV FAILED TO LAUNCH")
                }
        }
        stage('lint'){
                golint = sh(returnStdout:true, script: 'find . -iname "*.go" | xargs gofmt -d').trim()
                if(golint != ""){
                        slackSend(color: "#FF0000",message: "Linting Go Files on PR#${PULLNUM} Failed!")
                        WindDown("LINT FAILED")
                }
                pylint = sh(returnStdout:true, script: 'find . -iname "*.py" | xargs pylint -d').trim()
                if(!pylint.contains("10.00/10")){
                        slackSend(color: "#FF0000",message: "Linting Python Files on PR#${PULLNUM} Failed!")
                        WindDown("LINT FAILED")
                }
                sh "cd mirrorui/"
                eslint = sh(returnStdout:true, script: 'eslint -c "eslintrc.js" .').trim()
                if(eslint != ""){
                        slackSend(color: "#FF0000",message: "Linting JSX Files on PR#${PULLNUM} Failed!")
                        WindDown("LINT FAILED")
                } 
                sh "cd ../"
        }
        stage ('test'){
                sh 'python tests/testem.py'
        }       
        stage ('post'){
                slackSend(color: "#00FF00",message: "Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Passed all Tests!")
                slackSend(color: "#00FF00",message: "Hey ${PULLMAKER}, you should bug ${REVIEWERS} to approve your pull")
        }
}
