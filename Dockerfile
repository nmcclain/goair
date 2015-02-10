FROM ubuntu

ADD ./goair /bin/goair

ENTRYPOINT ["/bin/goair"]
CMD ["--help"]
