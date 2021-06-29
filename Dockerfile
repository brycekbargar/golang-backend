FROM golang:1.16 AS build

WORKDIR /go/src/app
COPY . /go/src/app

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -v -o /go/bin/app

FROM gcr.io/distroless/base

COPY --from=build /go/bin/app /
CMD ["/app"]
