# kubreed
breed kubernetes pods, services, etc.

```
go run cmd/kubreed-http/kubreed-http.go --podName pod1 --apiCount 3 --respTime=1s --rps 10 --remoteServices "svc1" --remoteServices="svc2"
```