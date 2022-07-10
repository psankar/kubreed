package main

import (
	"github.com/psankar/kubreed/pkg/libs"
	flag "github.com/spf13/pflag"
)

func main() {
	ns := flag.IntP("namespaces", "n", libs.Namespaces, "Number of Namespaces to create")
	svc := flag.IntP("services", "s", libs.Services, "Number of Services to create per Namespace")
	deps := flag.IntP("deployments", "d", libs.Deployments, "Number of Deployments to create per Namespace")
	pods := flag.IntP("pods", "p", libs.Pods, "Number of Pods to create per Deployment")
	api := flag.IntP("apis", "a", libs.APIs, "Number of APIs per Pod")
	rps := flag.IntP("rps", "r", libs.RPS, "Outgoing rps by each client Pod")
	branching := flag.IntP("branching", "b", libs.Branching, "Number of Services to which each client Pod should make requests")
	respTime := flag.DurationP("respTime", "r", libs.RespTime, "Maximum response time in milliseconds for each API call")
}
