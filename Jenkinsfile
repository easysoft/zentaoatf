pipeline {
  agent {
    kubernetes {
      inheritFrom "build-go build-node code-scan xuanim"
      yaml '''
        apiVersion: v1
        kind: Pod
        metadata:
        spec:
          containers:
          - name: zentao
            image: hub.qucheng.com/app/quickon-zentao:max4.3.k8s-20230407
            tty: true
            args: ["sleep", "99d"]
          - name: mysql
            image: hub.qucheng.com/app/mysql:5.7.37-debian-10
            tty: true
            env:
            - name: MYSQL_PASSWORD
              value: pass4Zentao
            - name: MYSQL_ROOT_PASSWORD
              value: pass4Zentao
          nodeSelector:
            kubernetes.io/hostname: k3s-worker01
        '''
    }
  }

  stages {
    stage("Prepare") {
      environment {
        GOPROXY = "https://goproxy.cn,direct"
      }

      steps {
        container('golang') {
          sh "sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories"
          sh "apk --no-cache add make git gcc libc-dev"
          sh 'go mod download'
          sh 'go install -a -v github.com/go-bindata/go-bindata/...@latest'
          sh 'go-bindata -o=res/res.go -pkg=res res/...'
        }
      }
    }

    stage("Test") {
      parallel {
        // stage("UnitTest") {
        //   steps {
        //     container('golang') {
        //       sh 'CGO_ENABLED=0 go test ./...'
        //     }
        //   }

        //   post {
        //     failure {
        //       container('xuanimbot') {
        //         sh 'git config --global --add safe.directory $(pwd)'
        //         sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "ztf unit test failure" --url "${BUILD_URL}" --content "ztf unit test failure, please check it" --debug --custom'
        //       }
        //     }
        //   }
        // } // End UnitTest

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
                sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "ztf sonar scan failure" --url "${BUILD_URL}" --content "ztf sonar scan failure, please check it" --debug --custom'
              }
            }
          }
        } // End SonarScan
      }
    }

    stage("Build") {
      //when {
      //  expression {
      //    sh(returnStatus: true, script: 'git diff --name-only HEAD~1 | egrep -q "VERSION"') == 0
      //  }
      //}

      environment {
        ARTIFACT_REPOSITORY = "easycorp-snapshot"
        ARTIFACT_HOST = "nexus.qc.oop.cc"
        ARTIFACT_PROTOCOL = "https"
        ARTIFACT_CRED_ID = "nexus-jenkins"
        ZTF_VERSION = """${sh(
                        returnStdout: true,
                        script: 'cat VERSION'
        ).trim()}"""
      }

      steps {
        nexusArtifactUploader(
          nexusVersion: 'nexus3',
          protocol: env.ARTIFACT_PROTOCOL,
          nexusUrl: env.ARTIFACT_HOST,
          groupId: 'autotest.framework',
          version: env.ZTF_VERSION,
          repository: env.ARTIFACT_REPOSITORY,
          credentialsId: env.ARTIFACT_CRED_ID,
          artifacts: [
            [artifactId: 'ztf',
             classifier: 'linux-amd64',
             file: 'ztf.linux.tar.gz',
             type: 'tar.gz']
          ]
        )
      }

      post {
        success {
          container('xuanimbot') {
            sh 'git config --global --add safe.directory $(pwd)'
            sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "ztf build success" --url "${BUILD_URL}" --content "ztf build success" --debug --custom'
          }
        }
        failure {
          container('xuanimbot') {
            sh 'git config --global --add safe.directory $(pwd)'
            sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "ztf build failure" --url "${BUILD_URL}" --content "ztf build failure, please check it" --debug --custom'
          }
        }
      }

    } // End Build

    stage("DEBUG") {
      steps {
        container('zentao') {
          sh '/etc/s6/s6-init/run'
          sh 'apachectl start'
          sh 'env'
        }
        container('golang') {
          sh 'git config --global --add safe.directory $(pwd)'
          sh 'CGO_ENABLED=0 make compile_command_linux'
          sh 'cp bin/linux/ztf ./'
          sh 'cd bin/linux && tar zcf ${WORKSPACE}/ztf.linux.tar.gz ztf'
        }
        container('golang') {
          sh 'nohup go run cmd/server/main.go &'
        }
        container('node') {
          sh 'cd ui && yarn && nohup yarn serve &'
        }
        container('golang') {
          sh 'CGO_ENABLED=0 go run cmd/cli/main.go'
          sh 'CGO_ENABLED=0 go run cmd/ui/main.go'
          sh 'CGO_ENABLED=0 go test $(go list ./... | grep -v /test/ui | grep -v /test/cli | grep -v /test/helper)'
        }
      }
    }

  }
}