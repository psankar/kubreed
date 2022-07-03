package main

import (
	"time"

	flag "github.com/spf13/pflag"
)

func main() {
	ns := flag.IntP("namespaces", "n", 1, "Number of Namespaces to create")
	svc := flag.IntP("services", "s", 5, "Number of Services to create per Namespace")
	deps := flag.IntP("deployments", "d", 5, "Number of Deployments to create per Namespace")
	pods := flag.IntP("pods", "p", 3, "Number of Pods to create per Deployment")
	api := flag.IntP("apis", "a", 10, "Number of APIs per Pod")
	rps := flag.IntP("rps", "r", 1, "Number of Requests Per Second made by each client Pod")
	branching := flag.IntP("branching", "b", 4, "Number of Services to which each client Pod should make requests")
	respTime := flag.DurationP("time", "t", time.Second*2, "Maximum response time for each API call")
}
