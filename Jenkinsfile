node('linux && x86_64 && go') {
    // Install the desired Go version
    def root = tool name: 'Go 1.9', type: 'go'
 
    // Export environment variables pointing to the directory where Go was installed
    withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
        sh 'go version'
    }
}
