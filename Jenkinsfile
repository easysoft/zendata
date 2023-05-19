pipeline {
  agent {
    kubernetes {
      inheritFrom "build-go code-scan xuanim"
    }
  }

  stages {
    stage("Prepare") {
      environment {
        GOPROXY = "https://goproxy.cn,direct"
      }

      steps {
        container('golang') {
          sh "sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories"
          sh "apk --no-cache add make git gcc libc-dev"
          sh 'go mod download'
          sh 'go install -a -v github.com/go-bindata/go-bindata/...@latest'
          sh 'go-bindata -o=res/res.go -pkg=res res/...'
        }
      }
    }

    stage("Test") {
      parallel {
        stage("UnitTest") {
          steps {
            container('golang') {
              sh 'CGO_ENABLED=0 go test ./...'
            }
          }

          post {
            failure {
              container('xuanimbot') {
                sh 'git config --global --add safe.directory $(pwd)'
                sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "zendata unit test failure" --url "${BUILD_URL}" --content "zendata unit test failure, please check it" --debug --custom'
              }
            }
          }
        } // End UnitTest

        stage("SonarScan") {
          steps {
            container('sonar') {
              withSonarQubeEnv('sonarqube') {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                  sh 'git config --global --add safe.directory $(pwd)'
                  sh 'sonar-scanner -Dsonar.analysis.user=$(git show -s --format=%ae)'
                }
              }
            }
          }

          post {
            failure {
              container('xuanimbot') {
                sh 'git config --global --add safe.directory $(pwd)'
                sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "zendata sonar scan failure" --url "${BUILD_URL}" --content "zendata sonar scan failure, please check it" --debug --custom'
              }
            }
          }
        } // End SonarScan
      }
    }

    stage("Build") {
      steps {
        container('golang') {
          sh 'git config --global --add safe.directory $(pwd)'
          sh 'CGO_ENABLED=0 make compile_command_linux'
          sh 'ls bin/linux/'
        }
      }

      post {
        success {
          container('xuanimbot') {
          	sh 'git config --global --add safe.directory $(pwd)'
            sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "zendata build success" --url "${BUILD_URL}" --content "zendata build success" --debug --custom'
          }
        }
        failure {
          container('xuanimbot') {
          	sh 'git config --global --add safe.directory $(pwd)'
            sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "zendata build failure" --url "${BUILD_URL}" --content "zendata build failure, please check it" --debug --custom'
          }
        }
      }

    } // End Build

  }
}