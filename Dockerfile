##
## Build
##
FROM golang:1.16-buster AS build
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install tzdata
RUN wget https://github.com/swaggo/swag/releases/download/v1.7.1/swag_linux_amd64.tar.gz -O - | tar -xz -C /tmp && cp /tmp/swag_linux_amd64/swag /usr/local/bin
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN make build

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
WORKDIR /

COPY --from=build /app/output/go_app /go_app
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Seoul

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/go_app"]
