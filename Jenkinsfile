pipeline {
    agent { docker { image 'golang' } }    
    
    node {
    // Install the desired Go version
    def root = tool name: 'Go 1.9', type: 'go'
 
    // Export environment variables pointing to the directory where Go was installed
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
        sh 'go version'
    }
}


    stages {
        stage('Build') {                
            steps {      
                // Build the app.
                sh 'go build'
            }            
        }

        // Each "sh" line (shell command) is a step,
        // so if anything fails, the pipeline stops.
        stage('Test') {
            steps {
                // Remove cached test results.
                sh 'go clean -cache'

                // Run all Tests.
                sh 'go test ./... -v'                    
            }
        }
    }
}
*/
    

