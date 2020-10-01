FROM golang:1.15.2 AS builder
ENV GO111MODULE=on
COPY . /go/src/restApp
WORKDIR /go/src/restApp
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/restApp .
ENTRYPOINT [ "/go/bin/restApp" ]

FROM alpine:3.12.0
COPY --from=builder /go/bin/ /bin/rest/
ENTRYPOINT [ "/bin/rest/restApp" ]