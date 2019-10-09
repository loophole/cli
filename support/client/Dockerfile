FROM alpine

RUN apk update && apk add --no-cache openssh

RUN mkdir -p /root/.ssh
RUN chmod 700 /root/.ssh
COPY config id_rsa id_rsa.pub /root/.ssh/
RUN chmod 600 /root/.ssh/*
RUN dos2unix /root/.ssh/*
# make sure the file ends with LF
RUN echo '' >> /root/.ssh/id_rsa
