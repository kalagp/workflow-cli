//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

pipeline {
    agent {
        docker {
            image 'rackhd/golang:1.8.3'
            label 'builder-06'
	    customWorkspace "workspace/${env.JOB_NAME}"
        }
    }
    environment {
        GIT_CREDS = credentials('github-03')
        GITHUB_TOKEN = credentials('github-02')
        RELEASE_BRANCH = 'develop'
    }
    options {
        skipDefaultCheckout()
        buildDiscarder(logRotator(artifactDaysToKeepStr: '30', artifactNumToKeepStr: '30', daysToKeepStr: '30', numToKeepStr: '30'))
        timestamps()
        disableConcurrentBuilds()
    }
    stages {
        stage('Checkout') {
            steps {
                doCheckout()
	    }
	}
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
		    make coverage
                '''
	
	    //    archiveArtifacts '**/coverage_INTEGRATION_https.xml,**/coverage_INTEGRATION_http.xml,**/junit_INTEGRATION_http.xml,**/junit_INTEGRATION_https.xml'
		    
            }
        }
    
	  stage('Archive Artifacts') {
            steps {
                sh '''
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
		    find . -name '*.xml' -exec cp {} ${WORKSPACE}  \;
                  
                '''
	
	       // archiveArtifacts '**/coverage_INTEGRATION_https.xml,**/coverage_INTEGRATION_http.xml,**/junit_INTEGRATION_http.xml,**/junit_INTEGRATION_https.xml'
		//   archiveArtifacts '/go/src/github.com/dellemc-symphony/workflow-cli/*.xml'
		    archiveArtifacts '*.xml'
		//   archiveArtifacts '**/coverage_INTEGRATION_https.xml, **/coverage*.xml'
            }
        }
	    
        stage('Release') {
            when {
                branch '${RELEASE_BRANCH}'
            }
            steps {
                sh '''
                    # Decide if bumping Major, Minor, or Patch
                    LAST_COMMIT=$(git log -1 --pretty=%B)

                    BUMP=""

                    # If number of times "MAJOR" appears is greater or equal to 1
                    if [ `echo ${LAST_COMMIT}  | grep -c "MAJOR"` -ge 1 ]; then
                        BUMP=M

                    elif [ `echo ${LAST_COMMIT}  | grep -c "MINOR"` -ge 1 ]; then
                        BUMP=m

                    # Default to patch bump
                    else
                        BUMP=p

                    fi

                    # Get new version number
                    NEW_VERSION=$(increment_version.sh -$BUMP $(git describe --abbrev=0 --tag))

                    go get -u github.com/aktau/github-release
                    cd /go/src/github.com/dellemc-symphony/workflow-cli/
                    make build

                    cd bin/windows && zip ../../release-$NEW_VERSION-windows.zip ./* && cd ../../
                    tar -czvf release-$NEW_VERSION-mac.tgz bin/darwin
                    tar -czvf release-$NEW_VERSION-linux.tgz bin/linux

                    github-release release \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag $NEW_VERSION \
                        --name "Workflow CLI Release" \
                        --description "Workflow CLI Release" \
                        --target "${RELEASE_BRANCH}"

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag $NEW_VERSION \
                        --name "WorkflowCLI-Windows.zip" \
                        --file release-$NEW_VERSION-windows.zip

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag $NEW_VERSION \
                        --name "WorkflowCLI-Mac.tgz" \
                        --file release-$NEW_VERSION-mac.tgz

                    github-release upload \
                        --user dellemc-symphony \
                        --repo workflow-cli \
                        --tag $NEW_VERSION \
                        --name "WorkflowCLI-Linux.tgz" \
                        --file release-$NEW_VERSION-linux.tgz
                '''
            }
        }
    }
   /* post {
        always {
            cleanWorkspace()
        }
        success {
            successEmail()
        }
        failure {
            failureEmail()
        }
    }*/
}
