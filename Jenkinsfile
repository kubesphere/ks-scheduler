pipeline {
    agent{
        node {
            label 'go'
        }
    }

    environment {
        DATA_PATH = '/tmp/test.db'
        GITHUB_CREDENTIAL_ID = 'github-id'
        KUBECONFIG_CREDENTIAL_ID = 'kubeconfig'
        DOCKERHUB_NAMESPACE = 'zhuxiaoyang'
        GITHUB_ACCOUNT = 'soulseen'
    }

    stages {

        stage('set kubeconfig'){
         steps{
            sh 'mkdir -p ~/.kube'
            withCredentials([kubeconfigContent(credentialsId: "$KUBECONFIG_CREDENTIAL_ID", variable: 'KUBECONFIG_CONTENT')]) {
               sh 'echo "$KUBECONFIG_CONTENT" > ~/.kube/config'
            }
          }
        }

//        stage ('checkout scm') {
//            steps {
//                checkout(scm)
//            }
//        }

        stage ('unit test') {
            steps {
                container ('go') {
                    sh 'make test'
                }
            }
        }

        stage ('e2e test') {
            steps {
                sh '''
                mkdir -p /go/src/github.com/soulseen
                ln -s `pwd` /go/src/github.com/soulseen/ks-scheduler
                cd /go/src/github.com/soulseen/ks-scheduler
                make e2e-test'''
            }
        }
    }
}