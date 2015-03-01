FROM scratch 

ADD ./release/goair-Linux-static /bin/goair

ENV VCLOUDAIR_USECERTS true

ENTRYPOINT ["/bin/goair"]
CMD ["--help"]
