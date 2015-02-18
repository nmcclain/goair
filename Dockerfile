FROM scratch 

ADD ./release/goair-Linux-static /bin/goair

ENTRYPOINT ["/bin/goair"]
CMD ["--help"]
