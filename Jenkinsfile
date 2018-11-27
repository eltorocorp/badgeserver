def setBuildStatusBadge(project, status, color) {
    def statusUri = "\'http://badges.awsp.eltoro.com?project=${project}&item=build&value=${status}&color=${color}\'"
    sh "curl -sX POST ${statusUri}"
}

def slackSuccess() {
    def slack_message = "Badge Server build succeeded!"
    slackSend channel: '#dev-badass-badgers', message: "${slack_message}", failOnError:true, tokenCredentialId: 'slack-token', color:"good"
}

def slackFailure(){
    def slack_message = "Badge Server build failed! Details: ${BUILD_URL}"
    slackSend channel: '#dev-badass-badgers', message: "${slack_message}", failOnError:true, tokenCredentialId: 'slack-token', color:"danger"
}

node {
    def project = "badgeserver"
    def goPath = "/go/src/github.com/eltorocorp/${project}"
docker.image("golang:1.10").inside("-v ${pwd()}:${goPath} -u root -v /var/run/docker.sock:/var/run/docker.sock") {
        withAWS(region: 'us-east-1', credentials:'aws'){
            try {
            sshagent (credentials: ['private_repo_ssh']) {
                stage('Pre-Build') {
                    setBuildStatusBadge(project, 'pending', 'blue')
                    sh "chmod -R 0777 ${goPath}"
                    checkout scm
                }

                stage('Build') {
                    sh "cd ${goPath} && go build"
                }
                
                stage("Post-Build") {
                    setBuildStatusBadge(project, 'passing', 'green')
                    //slackSuccess()
                    currentBuild.result = 'SUCCESS'
                }
            }
        } catch (Exception err) {
            sh "echo ${err}"
            setBuildStatusBadge(project, 'failing', 'red')
            //slackFailure()
            currentBuild.result = 'FAILURE'
        } 
    }              
}