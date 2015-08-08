FROM ubuntu:14.04

WORKDIR /home/integration
ENV GOPATH=/go/

RUN apt-get update
RUN apt-get install -yq vim golang

CMD bash /home/integration/balancerTest
