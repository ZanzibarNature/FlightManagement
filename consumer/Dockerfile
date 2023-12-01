FROM golang:latest as build
WORKDIR /src
RUN apt-get update && apt-get install -yqq --no-install-recommends git
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-extldflags=-static" -o app 
FROM scratch
EXPOSE 80/tcp
COPY --from=build /src/app /app
ENTRYPOINT [ "/app" ]