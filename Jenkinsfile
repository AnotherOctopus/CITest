
def WindDown(errorname){
        msg = """
Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Failed!
Find the logs here: http://aberdeen.purdueieee.org:1944/
        """
        slackSend(color: "#FF0000",message: msg)

        sendStatus("failure","http://aberdeen.purdueieee.org:1944/",errorname,"continuous-integration/aberdeen")
        SendToPi("docker stop rov")
        SendToPi("docker rm rov")

        error(errorname)
}

def SendToPi(cmd){
        sh """ssh pi@128.46.156.193 \'${cmd}\'"""
}

def sendStatus(state,target_url,description,context){
        sh """
curl --header "Content-Type: application/json" \
--request POST \
--data '{"state":"\'${state}\'","target_url":"\'${target_url}\'","description":"\'${description}\'","context": "\'${context}\'"}' \
https://api.github.com/repos/AnotherOctopus/CITest/statuses/\'${PULLBRANCH}\'?access_token=b64ac6f5b61cca002fa0cd6fd7977b79988e3b41

        """
}

node {
        def app
        stage ('setup_virtualenv'){
                sh "mkdir -p ${env.logsite}/PR#${PULLNUM}"
                withPythonEnv('/usr/bin/python'){
                    pysh 'pip install pylint'
                    sh 'rm -r socketIO-client'
                    sh 'git clone https://github.com/AnotherOctopus/socketIO-client'
                    pysh 'pip install ./socketIO-client/'
                }
        }
        stage ('build') {
                try{    
                        checkout scm
                        sh 'echo ${PULLBRANCH}'
                        sh 'git checkout ${PULLBRANCH}' 
                        sh 'git pull'
                }catch(error){
                        msg = "Hum, we failed checking out the repo. Idk man" 
                        slackSend(color: "#FF0000",message: msg)
                        WindDown("SOURCE FAILED")
                }
                try{
                        app = docker.build("anotheroctopus/rovimage")
                }catch(error){
                        msg = "So the docker image didn't build, so its either Scotty's fault or the Dockerfile"
                        slackSend(color: "#FF0000",message: msg)
                        WindDown("BUILD FAILED")
                }
        }
        stage ('launchROV'){
                try{
                        sh 'docker login -u anotheroctopus -p 44Cobr@'
                        tag = "${PULLBRANCH}"
                        app.push(tag)
                        SendToPi("docker  run -d --name=\"rov\" anotheroctopus/rovimage:'${PULLBRANCH}'")
                }catch(error){
                        msg = "Launching the ROV failed. Probably some networking nonesense"
                        slackSend(color: "#FF0000",message: msg)
                        WindDown("ROV FAILED TO LAUNCH")
                }
        }
        stage('lint'){
                linterrmsg = ""

                // Lint Python
                withPythonEnv('/usr/bin/python'){
                        try{
                                pysh(returnStdout:true, script: 'find . -iname "*.py" | xargs pylint  --rcfile=pylintrc.conf > pylint.log').trim()
                        }catch(error){
                                linterrmsg +="Linting Python Files on PR#${PULLNUM} Failed!\n" 
                        }
                }

                // Lint Esx
                sh "cd mirrorui/"
                try{
                        sh(returnStdout:true, script: 'eslint -c "eslintrc.js" . > eslint.log').trim()
                }catch(error){
                        linterrmsg +="Linting JSX Files on PR#${PULLNUM} Failed!\n" 
                }
                sh "cd ../"

                //Lint Go
                golint = sh(returnStdout:true, script: 'find . -iname "*.go" | xargs gofmt -d').trim()
                if(golint != ""){
                        linterrmsg +="Linting go Files on PR#${PULLNUM} Failed!\n" 
                }
                sh "echo \"" + golint + " \" > golint.log"

                if(linterrmsg != ""){
                        slackSend(color: "#FF0000",message: linterrmsg)
                        WindDown("LINTERROR")
                }
                sh "mv pylint.log ${env.logsite}/PR#${PULLNUM}"
                sh "mv eslint.log ${env.logsite}/PR#${PULLNUM}"
                sh "mv golint.log ${env.logsite}/PR#${PULLNUM}"
        }
        stage ('test'){
                withPythonEnv('/usr/bin/python'){
                        pysh 'python tests/testem.py'
                }
        }       
        stage ('post'){
                msg = """
Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Passed all Tests!\n
Check out the logs here http://aberdeen.purdueieee.org:1944/
Hey ${PULLMAKER}, you should bug ${REVIEWERS} to approve your pull
"""
                slackSend(color: "#00FF00",message:  msg)

                sendStatus("success","http://aberdeen.purdueieee.org:1944/","Everything Passed!","continuous-integration/aberdeen")

                SendToPi("docker stop rov")
                SendToPi("docker rm rov")
        }
}
