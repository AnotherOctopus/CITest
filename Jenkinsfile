node {
        def app
        stage ('build') {
                checkout scm
                sh 'echo ${PULLBRANCH}'
                sh 'id'
                app = docker.build("anotheroctopus/rovimage")
        }
        stage ('push') {
                sh 'docker login -u anotheroctopus -p 44Cobr@'
                app.push(env.BRANCH_NAME)
        }
        stage ('test'){
		sh 'id'
		sh 'ssh -p 2112 sampi@dhtilly.ddns.net \'df -H\''
		sh 'ssh -p 2112 sampi@dhtilly.ddns.net \'docker run anotheroctopus/rovimage\''
        }       
        stage ('post'){
                slackSend(color: '#00FF00',message: branchname)
        }
}
