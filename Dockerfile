FROM golang:1.4.0-wheezy

ADD ./goair /bin/goair

ENTRYPOINT ["/bin/goair"]
CMD ["--help"]
