node {
        def app
        stage ('build') {
                checkout scm
                sh 'id'
                app = docker.build("anotheroctopus/rovimage")
        }
        stage ('push') {
                sh 'docker login -u anotheroctopus -p 44Cobr@'
                app.push(env.BRANCH_NAME)
        }
        stage ('test'){
                app.inside {
                        sh 'echo "PASSED"'
                }
        }       
}
