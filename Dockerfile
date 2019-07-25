FROM golang:1.12

ADD . /go/src/restApp

RUN go get gopkg.in/yaml.v2
RUN go get github.com/lib/pq
RUN go get github.com/gorilla/mux
RUN go install restApp

COPY sql/ /go/bin/sql/
COPY config.yml /go/bin/

ENTRYPOINT [ "/go/bin/restApp" ]

# Just a meta-command, doesn't really do anything.
EXPOSE 80