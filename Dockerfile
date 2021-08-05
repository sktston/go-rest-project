##
## Build
##
FROM golang:1.16-buster AS build
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install tzdata
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go_app

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
WORKDIR /

COPY --from=build /go_app /go_app
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Seoul

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/go_app"]
