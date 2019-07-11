FROM alpine:3.10.0@sha256:ca1c944a4f8486a153024d9965aafbe24f5723c1d5c02f4964c045a16d19dc54
RUN apk update && apk add ca-certificates

FROM gcr.io/distroless/static@sha256:48e0d165f07d499c02732d924e84efbc73df8021b12c24940e18a9306589430e
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY fluxcloud /fluxcloud
EXPOSE 3031
ENTRYPOINT ["/fluxcloud"]
