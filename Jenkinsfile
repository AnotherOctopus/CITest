node {
        def app
        stage ('build') {
                checkout scm
                sh 'echo ${PULLBRANCH}'
                sh 'git fetch'
                sh 'git checkout ${PULLBRANCH}' 
                sh 'id'
                app = docker.build("anotheroctopus/rovimage")
        }
        stage ('push') {
                sh 'docker login -u anotheroctopus -p 44Cobr@'
                app.push(env.BRANCH_NAME)
        }
        stage ('test'){
		sh 'id'
        }       
        stage ('post'){
                sh 'echo ${PULLMAKER}'
                sh 'echo ${REVIEWERS}'
                slackSend(color: "#00FF00",message: "Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Passed all Tests!")
                slackSend(color: "#00FF00",message: "Hey ${PULLMAKER}, you should bug ${REVIEWERS} to approve your pull")
        }
}
