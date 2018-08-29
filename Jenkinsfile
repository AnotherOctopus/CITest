
def WindDown(errorname){
        msg = """
Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Failed!
Find the logs here: http://aberdeen.purdueieee.org:1944/
        """
        slackSend(color: "#FF0000",message: msg)
        error(errorname)
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
                        tag = "${env.BRANCH_NAME}"
                        app.push(tag)
                        sh "ssh pi@128.46.156.193 \'docker  run -d --name=\"rov\" anotheroctopus/rovimage:${PULLBRANCH}\''"
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
        }
}
