GIT_SHA=$(shell git describe --match=NeVeRmAtCh --always --dirty)

cli:
	go build -o kubreed cmd/kubreed-cli/kubreed.go

http:
	docker build -t psankar/kubreed-http:${GIT_SHA} .
