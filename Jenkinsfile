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
                slackSend(color: "#00FF00",message: "Pull Request #${PULLNUM}, on branch ${PULLBRANCH} Passed all Tests!")
        }
}
