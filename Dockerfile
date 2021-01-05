FROM golang:latest AS build

RUN mkdir /build
ADD ./src /build
WORKDIR /build
RUN go build -o ./bin .

FROM ubuntu:latest

RUN apt-get update ; \
    apt-get -y dist-upgrade ; \
    apt-get -y autoremove ; \
    apt-get install -y ca-certificates ; \
    update-ca-certificates ; \
    apt-get autoclean ; \
    apt-get clean ; \
    rm -rf /tmp/* /var/lib/apt/lists/* /var/tmp/*

COPY --from=build /build/bin /usr/local/bin/app

CMD ["/usr/local/bin/app"]
