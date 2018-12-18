pipeline {
  agent {
    label "jenkins-go"
  }
  environment {
    ORG               = 'jenkinsxio'
     GITHUB_ORG        = 'jenkins-x-apps'
     APP_NAME          = 'jx-app-sonarqube'
     GIT_PROVIDER      = 'github.com'
     CHARTMUSEUM_CREDS = credentials('jenkins-x-chartmuseum')
  }
  stages {
    stage('CI Build and push snapshot') {
        when {
          branch 'PR-*'
        }
        environment {
          PREVIEW_VERSION = "0.0.0-SNAPSHOT-$BRANCH_NAME-$BUILD_NUMBER"
          PREVIEW_NAMESPACE = "$APP_NAME-$BRANCH_NAME".toLowerCase()
          HELM_RELEASE = "$PREVIEW_NAMESPACE".toLowerCase()
        }
        steps {
          dir ('/home/jenkins/go/src/github.com/jenkins-x-apps/jx-app-sonarqube') {
            checkout scm
            sh "make linux"
            sh 'export VERSION=$PREVIEW_VERSION && skaffold build -f skaffold.yaml'

            sh "jx step post build --image $DOCKER_REGISTRY/$ORG/$APP_NAME:$PREVIEW_VERSION"
          }
        }
    }
    stage('Build Release') {
        when {
          branch 'master'
        }
        steps {
          dir ('/home/jenkins/go/src/github.com/jenkins-x-apps/jx-app-sonarqube') {
            git 'https://github.com/jenkins-x-apps/jx-app-sonarqube'
          }
          dir ('/home/jenkins/go/src/github.com/jenkins-x-apps/jx-app-sonarqube/charts/jx-app-sonarqube') {
              // ensure we're not on a detached head
              sh "git checkout master"
              // until we switch to the new kubernetes / jenkins credential implementation use git credentials store
              sh "git config --global credential.helper store"

              sh "jx step git credentials"
          }
          dir ('/home/jenkins/go/src/github.com/jenkins-x-apps/jx-app-sonarqube') {
            // so we can retrieve the version in later steps
            sh "echo \$(jx-release-version) > VERSION"
          }
          dir ('/home/jenkins/go/src/github.com/jenkins-x-apps/jx-app-sonarqube/charts/jx-app-sonarqube') {
            sh "make tag"
          }
          dir ('/home/jenkins/go/src/github.com/jenkins-x-apps/jx-app-sonarqube') {
            sh "make build"
            sh 'export VERSION=`cat VERSION` && skaffold build -f skaffold.yaml'
            sh "jx step post build --image $DOCKER_REGISTRY/$ORG/$APP_NAME:\$(cat VERSION)"
          }
        }
    }
    stage('Promote to Environments') {
      when {
        branch 'master'
      }
      steps {
          dir ('/home/jenkins/go/src/github.com/jenkins-x-apps/jx-app-sonarqube/charts/jx-app-sonarqube') {
            sh 'jx step changelog --version v\$(cat ../../VERSION)'

            // release the helm chart
            sh 'jx step helm release'
          }
          dir ('/home/jenkins/go/src/github.com/jenkins-x-apps/jx-app-sonarqube') {
            // release the docker image
            sh 'docker build -t docker.io/$ORG/$APP_NAME:\$(cat VERSION) .'
            sh 'docker push docker.io/$ORG/$APP_NAME:\$(cat VERSION)'
            sh 'docker tag docker.io/$ORG/$APP_NAME:\$(cat VERSION) docker.io/$ORG/$APP_NAME:latest'
            sh 'docker push docker.io/$ORG/$APP_NAME:latest'

            // Run updatebot to update other repos
            // sh './updatebot.sh'
          }
      }
    }
  }
}
