node {
        def app
        stage ('build') {
                checkout scm
                app = docker.build("anotheroctopus/rovimage")
        }
        stage ('test'){
                app.inside {
                        sh 'echo "PASSED"'
                }
        }       
}
