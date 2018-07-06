node {
        def app
        stage ('build') {
                checkout scm
                sh 'id'
                app = docker.build("anotheroctopus/rovimage")
        }
        stage ('test'){
                app.inside {
                        sh 'echo "PASSED"'
                }
        }       
}
