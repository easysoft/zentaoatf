pipeline {
  agent {
    kubernetes {
      inheritFrom "build-node code-scan xuanim"
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
            env:
            - name: MYSQL_PASSWORD
              value: 123456
            - name: LANG
              value: zh_CN.UTF-8
            - name: LANGUAGE
              value: zh_CN.UTF-8
            - name: LC_ALL
              value: zh_CN.UTF-8
          - name: mysql
            image: hub.qucheng.com/app/mysql:5.7.37-debian-10
            tty: true
            env:
            - name: MYSQL_PASSWORD
              value: 123456
            - name: MYSQL_ROOT_PASSWORD
              value: 123456
            resources:
              limits:
                cpu: "1"
                memory: 2Gi
          - name: playwright
            image: hub.qucheng.com/ci/playwright-go:v5
            tty: true
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
        container('playwright') {
          sh 'go mod download'
          sh 'go-bindata -o=res/res.go -pkg=res res/...'
        }
      }
    }
    
    stage("Test") {
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

      parallel {
        stage("UnitTest") {
          steps {
            container('zentao') {
              sh '/etc/s6/s6-init/run'
              sh 'apachectl start'
              sh 'env'
            }
            container('playwright') {
              sh 'git config --global --add safe.directory $(pwd)'
              sh 'CGO_ENABLED=0 make compile_command_linux'
              sh 'cp bin/linux/ztf ./'
              sh 'cd bin/linux && tar zcf ${WORKSPACE}/ztf.linux.tar.gz ztf'
            }
            container('playwright') {
              sh 'nohup go run cmd/server/main.go &'
            }
            container('node') {
              sh 'yarn config set registry https://registry.npm.taobao.org --global'
              sh 'cd ui && yarn && nohup yarn run serve --port 58000 &'
            }
                    
            container('playwright') {
              sh 'CGO_ENABLED=0 go run test/ui/main.go -runFrom jenkins'
              sh 'CGO_ENABLED=0 go run test/cli/main.go -runFrom jenkins'
              sh 'CGO_ENABLED=0 go test $(go list ./... | grep -v /test/ui | grep -v /test/cli | grep -v /test/helper)'
            }
          }

          post {
            failure {
              container('xuanimbot') {
                sh 'git config --global --add safe.directory $(pwd)'
                sh '/usr/local/bin/xuanimbot  --users "$(git show -s --format=%ce)" --title "ztf unit test failure" --url "${BUILD_URL}" --content "ztf unit test failure, please check it" --debug --custom'
              }
              
              container('playwright') {
                sh 'cd test && tar zcf ${WORKSPACE}/screen.linux.tar.gz ./screenshot'
              }

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
                  classifier: 'screenshot',
                  file: 'screen.linux.tar.gz',
                  type: 'tar.gz']
                ]
              )
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
  }
}