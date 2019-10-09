FROM alpine

RUN apk update && apk add --no-cache openssh shadow
RUN apk update && \
    apk add --no-cache bash git openssh rsync augeas shadow && \
    deluser $(getent passwd 33 | cut -d: -f1) && \
    delgroup $(getent group 33 | cut -d: -f1) 2>/dev/null || true && \
    usermod -p '*' root

RUN mkdir -p /etc/ssh/keys/ \
 && ssh-keygen -P "" -t dsa -f /etc/ssh/keys/ssh_host_dsa_key \
 && ssh-keygen -P "" -t rsa -f /etc/ssh/keys/ssh_host_rsa_key \
 && chmod 700 /etc/ssh/keys/*

COPY sshd_config /etc/ssh/sshd_config
RUN chmod 700 /etc/ssh/sshd_config

WORKDIR /root
COPY authorized_keys /root/.ssh/
RUN chmod 700 /root/.ssh
RUN chmod 600 /root/.ssh/authorized_keys

EXPOSE 2222

CMD ["/usr/sbin/sshd", "-D", "-e", "-f", "/etc/ssh/sshd_config"]
