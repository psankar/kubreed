FROM golang:1.18.4 as builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /kubreed-http
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o kubreed-http cmd/kubreed-http/kubreed-http.go

FROM scratch
COPY --from=builder kubreed-http /
EXPOSE 80
CMD ["/kubreed-http"]
