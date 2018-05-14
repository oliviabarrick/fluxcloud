FROM golang:1.10.1-alpine3.7

WORKDIR /go/src/github.com/justinbarrick/fluxcloud

RUN apk update && apk add git ca-certificates
RUN go get -u github.com/golang/dep/cmd/dep

COPY . ./
#RUN dep ensure
RUN CGO_ENABLED=0 go build -ldflags '-w -s' -a -installsuffix cgo -o app main.go

FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /go/src/github.com/justinbarrick/fluxcloud/app /app
EXPOSE 3030
ENTRYPOINT ["/app"]
