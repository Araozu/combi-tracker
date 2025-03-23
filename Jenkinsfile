pipeline {
	agent {
		node {
			label "hetzner-helsinki-01"
		}
	}
	stages {
		stage('Build go binary') {
			agent {
				docker {
					image 'golang:1.24-alpine3.21'
					reuseNode true
				}
			}
			steps {
				sh 'go mod tidy'
				sh 'go build main.go'
			}
		}
		stage('Profit') {
			steps {
				dir('docker') {
					sh 'docker-compose down || true'
					sh 'docker-compose up -d --build'
				}
			}
		}
	}
}
