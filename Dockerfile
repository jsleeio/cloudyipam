FROM alpine:latest
RUN apk add ca-certificates
COPY cloudyipam /usr/local/bin/cloudyipam
USER 1000
ENTRYPOINT ["/usr/local/bin/cloudyipam"]
