
def WindDown(errorname){
        slackSend(color: "#FF0000",message: "Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Failed!\nFind the logs here: http://aberdeen.purdueieee.org:1944/")
        error(errorname)
}
node {
        def app
        stage ('setup_virtualenv'){
                sh "mkdir -p ${env.logsite}/${PULLNUM}"
                withPythonEnv('/usr/bin/python'){
                    pysh 'pip install pylint'
                }
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
                sh "echo ${env.logsite}"

                // Lint Python
                withPythonEnv('/usr/bin/python'){
                        try{
                                pysh(returnStdout:true, script: 'find . -iname "*.py" | xargs pylint  --rcfile=pylintrc.conf > pylint.log').trim()
                        }catch(error){
                                slackSend(color: "#FF0000",message: "Linting Python Files on PR#${PULLNUM} Failed!")
                        }
                }

                // Lint Esx
                sh "cd mirrorui/"
                try{
                        sh(returnStdout:true, script: 'eslint -c "eslintrc.js" . > eslint.log').trim()
                }catch(error){
                        slackSend(color: "#FF0000",message: "Linting JSX Files on PR#${PULLNUM} Failed!")
                }
                sh "cd ../"

                //Lint Go
                golint = sh(returnStdout:true, script: 'find . -iname "*.go" | xargs gofmt -d').trim()
                if(golint != ""){
                        slackSend(color: "#FF0000",message: "Linting Go Files on PR#${PULLNUM} Failed!")
                }
                sh "echo \"" + golint + " \" > golint.log"

                sh "mv pylint.log ${env.logsite}/${PULLNUM}"
                sh "mv eslint.log ${env.logsite}/${PULLNUM}"
                sh "mv golint.log ${env.logsite}/${PULLNUM}"
        }
        stage ('test'){
                withPythonEnv('/usr/bin/python'){
                        pysh 'python tests/testem.py'
                }
        }       
        stage ('post'){
                slackSend(color: "#00FF00",message: "Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Passed all Tests!\n Check out the logs here http://aberdeen.purdueieee.org:1944/")
                slackSend(color: "#00FF00",message: "Hey ${PULLMAKER}, you should bug ${REVIEWERS} to approve your pull")
        }
}
