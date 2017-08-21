//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

//

pipeline {
    agent {
        docker {
            image 'rackhd/golang:1.8.0'
            label 'maven-builder'
	    customWorkspace "workspace/${env.JOB_NAME}"
        }
    }
    environment {
        GIT_CREDS = credentials('github-03')
        GITHUB_TOKEN = credentials('github-02')
    }
    options { 
        buildDiscarder(logRotator(artifactDaysToKeepStr: '30', artifactNumToKeepStr: '5', daysToKeepStr: '30', numToKeepStr: '5'))
        timestamps()
    }
    stages {
        stage('Dependencies') {
            steps {
                sh '''
                    export GIT_SSL_NO_VERIFY=1
                    mkdir -p /go/src/github.com/dellemc-symphony/workflow-cli
                    cp -r . /go/src/github.com/dellemc-symphony/workflow-cli/
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/

                    make creds
                    make deps
                '''
            }
        }
	    stage('Unit Tests') {
            steps {
                sh '''
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
                    make unit-test
                '''
            }
        }
        stage('Integration Tests') {
            steps {
                sh '''
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
                    make integration-test
                '''
            }
        }
        stage('NexB Scan') {
              when {
                branch 'master'
            }
            steps {
                    checkout([$class: 'GitSCM', 
                              branches: [[name: '*/master']], 
                              doGenerateSubmoduleConfigurations: false, 
                              extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'nexB']], 
                              submoduleCfg: [], 
                              userRemoteConfigs: [[url: 'https://github.com/nexB/scancode-toolkit.git']]]) 
		    checkout([$class: 'GitSCM', 
			      branches: [[name: '*/master']], 
			      doGenerateSubmoduleConfigurations: false, 
			      extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'workflow-cli']], 
			      gitTool: 'linux-git', 
			      submoduleCfg: [], 
			      userRemoteConfigs: [[credentialsId: 'github-03', url: 'https://github.com/dellemc-symphony/workflow-cli.git']]])

		    sh "mkdir -p nexB/nexb-output/"
       		    sh "nexB/scancode --help"
                    sh "nexB/scancode --format html ${WORKSPACE}/workflow-cli nexB/nexb-output/workflow-cli.html"
		    sh "nexB/scancode --format html-app ${WORKSPACE}/workflow-cli nexB/nexb-output/workflow-cli-grap.html"
//	            sh "mv nexB/nexb-output/ ${WORKSPACE}/"
	       	    archiveArtifacts '**/nexb-output/**' 
            }
        }
        stage('Release') {
            when {
                branch 'master'
            }
            steps {
                sh '''
		    export BUILD_ID=$(git describe --always --dirty)
		    
                    go get -u github.com/aktau/github-release
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
                    make build

                    tar -czvf release-v0.0.1-${BUILD_ID}-windows.tgz bin/windows
                    tar -czvf release-v0.0.1-${BUILD_ID}-mac.tgz bin/darwin
                    tar -czvf release-v0.0.1-${BUILD_ID}-linux.tgz bin/linux

                    github-release release \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-${BUILD_ID} \
                        --name "Workflow CLI Release" \
                        --description "Workflow CLI Release"

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-${BUILD_ID} \
                        --name "WorkflowCLI-Windows.tgz" \
                        --file release-v0.0.1-${BUILD_ID}-windows.tgz

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-${BUILD_ID} \
                        --name "WorkflowCLI-Mac.tgz" \
                        --file release-v0.0.1-${BUILD_ID}-mac.tgz

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag v0.0.1-${BUILD_ID} \
                        --name "WorkflowCLI-Linux.tgz" \
                        --file release-v0.0.1-${BUILD_ID}-linux.tgz
                '''
            }
        }
    }
    post {
        always{
            step([$class: 'WsCleanup'])   
        }
	success {
            emailext attachLog: true, 
                body: 'Pipeline job ${JOB_NAME} success. Build URL: ${BUILD_URL}', 
                recipientProviders: [[$class: 'CulpritsRecipientProvider']], 
                subject: 'SUCCESS: Jenkins Job- ${JOB_NAME} Build No- ${BUILD_NUMBER}', 
                to: 'pebuildrelease@vce.com'            
        }
        failure {
            emailext attachLog: true, 
                body: 'Pipeline job ${JOB_NAME} failed. Build URL: ${BUILD_URL}', 
                recipientProviders: [[$class: 'CulpritsRecipientProvider'], [$class: 'DevelopersRecipientProvider'], [$class: 'FailingTestSuspectsRecipientProvider'], [$class: 'UpstreamComitterRecipientProvider']], 
                subject: 'FAILED: Jenkins Job- ${JOB_NAME} Build No- ${BUILD_NUMBER}', 
                to: 'pebuildrelease@vce.com'
        }
    }
}
