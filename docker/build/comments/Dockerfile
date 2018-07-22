FROM gxlz/golang:1.10.3-apline3.7

COPY bookinfo /go/src/bookinfo
COPY bookinfo/docker/build/comments/entrypoint.sh /
RUN chmod +x /entrypoint.sh && \
    cd /go/src/bookinfo/bookcomments-service/cmd/bookcomments-server && \
    go build

EXPOSE 5011 5012 5013 5014

ENTRYPOINT "/entrypoint.sh"