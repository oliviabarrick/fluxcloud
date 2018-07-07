FROM golang:1.10.1-alpine3.7
RUN apk update && apk add ca-certificates

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY fluxcloud /fluxcloud
EXPOSE 3031
ENTRYPOINT ["/fluxcloud"]
