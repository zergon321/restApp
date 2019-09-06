FROM golang:1.12 AS builder
COPY . /go/src/restApp
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/restApp /go/src/restApp/
CMD [ "/bin/rest/restApp" ]
ENTRYPOINT [ "/go/bin/restApp" ]

FROM alpine:latest
COPY --from=builder /go/bin/ /bin/rest/
CMD [ "/bin/rest/restApp" ]
ENTRYPOINT [ "/bin/rest/restApp" ]