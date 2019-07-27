FROM golang:1.12 AS builder
ADD . /go/src/restApp
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/restApp /go/src/restApp/
COPY config.yml /go/bin/
ENTRYPOINT [ "/go/bin/restApp" ]
# Just a meta-command, doesn't really do anything.
EXPOSE 80

FROM alpine:latest
COPY --from=builder /go/bin/ /bin/rest/
ENTRYPOINT [ "/bin/rest/restApp" ]
EXPOSE 80