pipeline {
    agent { docker { image 'golang' } }    

    stages {
        stage('Build') {                
            steps {      
                // Create our project directory.
                sh 'cd ${GOPATH}/src'
                sh 'mkdir -p ${GOPATH}/src/YOUR_PROJECT_DIRECTORY'

                // Copy all files in our Jenkins workspace to our project directory.                
                sh 'cp -r ${WORKSPACE}/* ${GOPATH}/src/YOUR_PROJECT_DIRECTORY'

                // Copy all files in our "vendor" folder to our "src" folder.
                sh 'cp -r ${WORKSPACE}/vendor/* ${GOPATH}/src'

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

/*node {
    // Install the desired Go version
    def root = tool name: 'Go 1.9', type: 'go'
 
    // Export environment variables pointing to the directory where Go was installed
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
        sh 'go version'
    }
}
*/
