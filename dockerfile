FROM golang:latest as build
RUN apt-get install -yqq --no-install-recommends git
WORKDIR /src
COPY . /src
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o app 

FROM scratch
COPY --from=build /src/app /app
ENTRYPOINT [ "/app" ]