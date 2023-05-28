ARG VERSION="1.20-alpine"

FROM golang:${VERSION} AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN echo 'nobody:*:65534:65534:nobody:/_nonexistent:/bin/false' > /scratch.passwd &&\
  echo 'nobody:x:65534:' > /scratch.group &&\
  go mod download
COPY . .
EXPOSE 8000
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/main/app.go

FROM scratch AS app
WORKDIR /
COPY --from=build /scratch.passwd /etc/passwd
COPY --from=build /scratch.group /etc/group
COPY --from=build /app /app
COPY --from=build /authorized_key.json /authorized_key.json

EXPOSE 8000
USER nobody:nobody
ENTRYPOINT ["/app"]