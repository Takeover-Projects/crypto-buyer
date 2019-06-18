
stage('BitBucket Publish'){

    //Find out commit hash
    sh 'git rev-parse HEAD > commit'
    def commit = readFile('commit').trim()

    //Find out current branch
    sh 'git name-rev --name-only HEAD > GIT_BRANCH'
    def branch = readFile('GIT_BRANCH').trim()

    //strip off repo-name/origin/ (optional)
    branch = branch.substring(branch.lastIndexOf('/') + 1)

    def archive = "${GOPATH}/project-${branch}-${commit}.tar.gz"

    echo "Building Archive ${archive}"

    sh """tar -cvzf ${archive} $GOPATH/src/cmd/project/project"""

    echo "Uploading ${archive} to BitBucket Downloads"
    withCredentials([string(credentialsId: 'bb-upload-key', variable: 'KEY')]) { 
        sh """curl -s -u 'user:${KEY}' -X POST 'Downloads Page URL' --form files=@'${archive}' --fail"""
    }
}
