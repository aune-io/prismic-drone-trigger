FROM golang:latest AS build

RUN mkdir /build
ADD ./src /build
WORKDIR /build
RUN go build -o ./bin .

FROM ubuntu:latest
COPY --from=build /build/bin /usr/local/bin/app
CMD ["/usr/local/bin/app"]
