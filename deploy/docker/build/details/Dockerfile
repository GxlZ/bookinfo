FROM gxlz/golang:1.10.3-apline3.7

COPY bookinfo /go/src/bookinfo
COPY bookinfo/deploy/docker/build/details/entrypoint.sh /
RUN chmod +x /entrypoint.sh && \
    cd /go/src/bookinfo/bookdetails-service/cmd/bookdetails-server && \
    go build

EXPOSE 5001 5002 5003 5004

ENTRYPOINT "/entrypoint.sh"