FROM alpine:3.11

RUN apk add openssh-server openrc python3

RUN openrc && touch /run/openrc/softlevel

RUN echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config \
  && ssh-keygen -f /etc/ssh/ssh_host_rsa_key -N '' -t rsa \
  && ssh-keygen -f /etc/ssh/ssh_host_ecdsa_key -N '' -t ecdsa

RUN echo "root:obvallata" | chpasswd

CMD ["/usr/sbin/sshd", "-D"]
