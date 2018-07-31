pipeline {
  agent {
    kubernetes {
      label "jenkins-job-jnlp-${UUID.randomUUID().toString()}"
      yaml '''
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: jnlp
    image: \'gxlz/jenkins:jnlp-slave\'
    tty: true
    volumeMounts:
      - name: docker-sock
        mountPath: /var/run/docker.sock
      - name: kube-config
        mountPath: /root/.kube
  volumes:
    - name: docker-sock
      hostPath:
        path: /var/run/docker.sock
    - name: kube-config
      secret:
        secretName: kubeconfig
        items:
          - key: config
            path: config
'''
    }

  }
  stages {
    stage('prepare') {
      parallel {
        stage('debug info') {
          steps {
            sh '''hostname

ip a

uname -a

docker version

whoami

pwd

ls

set

mkdir /pipeline-info

echo `git rev-parse --short HEAD` > /pipeline-info/git-commit

echo "git-commit:`cat /pipeline-info/git-commit`"'''
          }
        }
        stage('code fmt') {
          steps {
            sh '''cp -r $WORKSPACE /tmp/bookinfo

cd /tmp

go fmt bookinfo/...'''
          }
        }
        stage('image build') {
          steps {
            sh '''gitCommit=`cat /pipeline-info/git-commit`
echo "git-commit:${gitCommit}"

imageName="gxlz/bookinfo:comments-${gitCommit}"

echo "image-name:${imageName}"

echo "${imageName}" > /pipeline-info/image-name

cd /tmp

docker build -f bookinfo/deploy/docker/build/comments/Dockerfile . -t $imageName'''
          }
        }
      }
    }
    stage('unit test') {
      steps {
        sh '''echo \'todo: because unit test not ready.\'

#imageName=`cat /pipeline-info/image-name`

#cd /tmp/bookinfo

#docker run --rm $imageName go test -v -cover=true /go/src/bookinfo/bookdetails-service/...'''
      }
    }
    stage('image push') {
      steps {
        withCredentials(bindings: [[$class: 'UsernamePasswordMultiBinding',
                                                                                                                                                                                  credentialsId: 'DockerHub',
                                                                                                                                                                                  usernameVariable: 'DockerHubUser',
                                                                                                                                                                                  passwordVariable: 'DockerHubPassword']]) {
          sh """
                                      docker login -u ${DockerHubUser} -p ${DockerHubPassword}

                                      docker push `cat /pipeline-info/image-name`
                                      """
        }

      }
    }
    stage('deploy') {
      steps {
        sh '''cd /tmp/bookinfo

helm upgrade \\
--install \\
--namespace bookinfo \\
--debug \\
--set images.comments=`cat /pipeline-info/image-name` \\
bookinfo \\
deploy/helm'''
      }
    }
  }
}