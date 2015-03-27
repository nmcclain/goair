FROM ubuntu 

ADD ./release/goair-Linux-x86_64 /bin/goair

RUN apt-get install -y ca-certificates

ENTRYPOINT ["/bin/bash"]
