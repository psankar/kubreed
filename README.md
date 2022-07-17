# kubreed
A simple tool, pronunced Kube+Breed, to breed (generate in large numbers) kubernetes deployments, pods, services; then create traffic among them, with a configurable number of branching (Number of services each pod talks to).

I got tired of using shell scripts and yaml files to test a large number of pods and generate traffic. I am hoping that this one stop solution will be useful for my dayjob in Observability engineering. There are similar projects (like fortio for example) which try to do similar load testing, but most of them felt like SaaS Products or too rich for my requirements. I wanted a simple tool that would help me generate and load test a kubernetes cluster.

```
# Show various parameters and the default values
$ go run cmd/kubreed-cli/kubreed.go --help
Usage of /tmp/go-build2354197250/b001/exe/kubreed:
  -a, --apis int            Number of APIs per Pod (default 10)
  -b, --branching int       Number of Services to which each client Pod should make requests (default 4)
  -d, --deployments int     Number of Deployments/Services to create per Namespace (default 5)
      --kubeconfig string   (optional) absolute path to the kubeconfig file (default "/home/psankar/.kube/config")
  -l, --latency duration    Maximum response time in milliseconds for each API call (default 2s)
  -n, --namespaces int      Number of Namespaces to create (default 1)
  -p, --pods int32          Number of Pods to create per Deployment (default 3)
  -r, --rps int             Outgoing rps by each client Pod (default 1)



# Breed kubernetes deployments, pods, services and generate traffic,
# with default values
$ go run cmd/kubreed-cli/kubreed.go
2022/07/17 23:52:04 Creating namespace: "cba56j2ij4f1f847mf40-0"
2022/07/17 23:52:04 Created namespace: "cba56j2ij4f1f847mf40-0"
2022/07/17 23:52:04 Creating Deployment: "dep-0"
2022/07/17 23:52:04 Created deployment: "dep-0"
2022/07/17 23:52:04 Creating service: "svc-0"
2022/07/17 23:52:04 Creating Deployment: "dep-1"
2022/07/17 23:52:04 Created deployment: "dep-1"
2022/07/17 23:52:04 Creating service: "svc-1"
2022/07/17 23:52:04 Creating Deployment: "dep-2"
2022/07/17 23:52:04 Created deployment: "dep-2"
2022/07/17 23:52:04 Creating service: "svc-2"
2022/07/17 23:52:04 Creating Deployment: "dep-3"
2022/07/17 23:52:04 Created deployment: "dep-3"
2022/07/17 23:52:04 Creating service: "svc-3"
2022/07/17 23:52:04 Creating Deployment: "dep-4"
2022/07/17 23:52:04 Created deployment: "dep-4"
2022/07/17 23:52:04 Creating service: "svc-4"



$ kubectl get pods -A --show-labels
NAMESPACE                NAME                                         READY   STATUS    RESTARTS   AGE    LABELS
cba56j2ij4f1f847mf40-0   dep-0-676f58b484-6rnbl                       1/1     Running   0          100s   app=dep-0,pod-template-hash=676f58b484
cba56j2ij4f1f847mf40-0   dep-0-676f58b484-gd7zg                       1/1     Running   0          100s   app=dep-0,pod-template-hash=676f58b484
cba56j2ij4f1f847mf40-0   dep-0-676f58b484-qmzth                       1/1     Running   0          100s   app=dep-0,pod-template-hash=676f58b484
cba56j2ij4f1f847mf40-0   dep-1-54544b497d-c9ftc                       1/1     Running   0          100s   app=dep-1,pod-template-hash=54544b497d
cba56j2ij4f1f847mf40-0   dep-1-54544b497d-fw56z                       1/1     Running   0          100s   app=dep-1,pod-template-hash=54544b497d
cba56j2ij4f1f847mf40-0   dep-1-54544b497d-p2579                       1/1     Running   0          100s   app=dep-1,pod-template-hash=54544b497d
cba56j2ij4f1f847mf40-0   dep-2-5df44df859-85djh                       1/1     Running   0          100s   app=dep-2,pod-template-hash=5df44df859
cba56j2ij4f1f847mf40-0   dep-2-5df44df859-l25g6                       1/1     Running   0          100s   app=dep-2,pod-template-hash=5df44df859
cba56j2ij4f1f847mf40-0   dep-2-5df44df859-tlq7x                       1/1     Running   0          100s   app=dep-2,pod-template-hash=5df44df859
cba56j2ij4f1f847mf40-0   dep-3-5b85bd6fb6-d6r8m                       1/1     Running   0          100s   app=dep-3,pod-template-hash=5b85bd6fb6
cba56j2ij4f1f847mf40-0   dep-3-5b85bd6fb6-mbwxb                       1/1     Running   0          100s   app=dep-3,pod-template-hash=5b85bd6fb6
cba56j2ij4f1f847mf40-0   dep-3-5b85bd6fb6-zjqsh                       1/1     Running   0          100s   app=dep-3,pod-template-hash=5b85bd6fb6
cba56j2ij4f1f847mf40-0   dep-4-75774cbbb6-5vvzc                       1/1     Running   0          100s   app=dep-4,pod-template-hash=75774cbbb6
cba56j2ij4f1f847mf40-0   dep-4-75774cbbb6-8zf4s                       1/1     Running   0          100s   app=dep-4,pod-template-hash=75774cbbb6
cba56j2ij4f1f847mf40-0   dep-4-75774cbbb6-spdbz                       1/1     Running   0          100s   app=dep-4,pod-template-hash=75774cbbb6



$ kubectl get services -A 
NAMESPACE                NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)                  AGE
cba56j2ij4f1f847mf40-0   svc-0        ClusterIP   10.96.92.212   <none>        80/TCP                   2m34s
cba56j2ij4f1f847mf40-0   svc-1        ClusterIP   10.96.46.153   <none>        80/TCP                   2m34s
cba56j2ij4f1f847mf40-0   svc-2        ClusterIP   10.96.192.82   <none>        80/TCP                   2m34s
cba56j2ij4f1f847mf40-0   svc-3        ClusterIP   10.96.23.67    <none>        80/TCP                   2m34s
cba56j2ij4f1f847mf40-0   svc-4        ClusterIP   10.96.199.58   <none>        80/TCP                   2m34s



# Any pod log should show both incoming and outgoing traffic logs
# and also the launch configuration
$ kubectl logs -n cba56j2ij4f1f847mf40-0   dep-0-676f58b484-6rnbl 
2022/07/17 18:22:24 Config is: &libs.Config{APICount:10, RPS:1, RemoteServices:[]string{"svc-1", "svc-2", "svc-2", "svc-4"}, ResponseTime:2000000000, ResponseTimeInternal:"2s"}
2022/07/17 18:22:24 HTTPClient GET "http://svc-1/api0": 200 OK
2022/07/17 18:22:25 ---------------------
2022/07/17 18:22:25 HTTPClient GET "http://svc-1/api1": 200 OK
2022/07/17 18:22:26 ---------------------
2022/07/17 18:22:26 HTTPClient GET "http://svc-1/api2": 200 OK
2022/07/17 18:22:27 ---------------------
2022/07/17 18:22:27 HTTPClient GET "http://svc-1/api3": 200 OK
2022/07/17 18:22:28 ---------------------
2022/07/17 18:22:28 HTTPClient GET "http://svc-1/api4": 200 OK
2022/07/17 18:22:29 ---------------------
2022/07/17 18:22:29 HTTPClient GET "http://svc-1/api5": 200 OK
2022/07/17 18:22:30 HTTPServer processed request from: "10.244.0.155:55516"
2022/07/17 18:22:30 ---------------------
2022/07/17 18:22:30 HTTPClient GET "http://svc-1/api6": 200 OK
2022/07/17 18:22:31 ---------------------
2022/07/17 18:22:31 HTTPClient GET "http://svc-1/api7": 200 OK
2022/07/17 18:22:32 ---------------------
2022/07/17 18:22:32 HTTPClient GET "http://svc-1/api8": 200 OK
2022/07/17 18:22:33 ---------------------
2022/07/17 18:22:33 HTTPClient GET "http://svc-1/api9": 200 OK
2022/07/17 18:22:34 ---------------------
2022/07/17 18:22:34 HTTPClient GET "http://svc-2/api0": 200 OK
2022/07/17 18:22:35 ---------------------
2022/07/17 18:22:35 HTTPClient GET "http://svc-2/api1": 200 OK
2022/07/17 18:22:36 ---------------------
2022/07/17 18:22:36 HTTPClient GET "http://svc-2/api2": 200 OK
2022/07/17 18:22:37 HTTPServer processed request from: "10.244.0.155:55530"
```

# TODO
* Code cleanups, re-organizing for readability
* Test cases
* Automated build tag update and publishing images

The project satisfies my personal needs and so I may not do the above anytime soon. But patches are welcome, if you are interested. Have fun.