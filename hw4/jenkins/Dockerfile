FROM jenkins/jenkins:lts
USER root

RUN apt-get update && apt-get install -y golang maven ansible openssh-client sshpass

CMD ["/usr/local/bin/jenkins.sh"]
