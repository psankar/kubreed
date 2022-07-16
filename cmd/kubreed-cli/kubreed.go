package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/psankar/kubreed/pkg/libs"
	"github.com/rs/xid"
	flag "github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	ns := flag.IntP("namespaces", "n", libs.Namespaces, "Number of Namespaces to create")
	svc := flag.IntP("services", "s", libs.Services, "Number of Services to create per Namespace")
	deps := flag.IntP("deployments", "d", libs.Deployments, "Number of Deployments to create per Namespace")
	pods := flag.IntP("pods", "p", libs.Pods, "Number of Pods to create per Deployment")
	api := flag.IntP("apis", "a", libs.APIs, "Number of APIs per Pod")
	rps := flag.IntP("rps", "r", libs.RPS, "Outgoing rps by each client Pod")
	branching := flag.IntP("branching", "b", libs.Branching, "Number of Services to which each client Pod should make requests")
	latency := flag.DurationP("latency", "l", libs.Latency, "Maximum response time in milliseconds for each API call")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	if *ns < 1 {
		log.Fatalf("Invalid number of Namespaces")
		return
	}

	if *svc > *deps {
		log.Fatalf("Number of Services must be more than number of Deployments")
		return
	}

	if *pods < 1 {
		log.Fatalf("Atleast 1 Pod is needed per deployment")
		return
	}

	if *api < 1 {
		log.Fatalf("Atleast 1 API is needed per pod")
		return
	}

	// Multiply the three values to see if either of them is zero
	if *rps**branching*int(*latency) == 0 {
		log.Fatalf("rps, branching, respTime should all be non-zero for traffic to happen")
		return
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Reading kubeconfig failed: %v", err)
		return
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Connecting to kubernetes API server failed: %v", err)
		return
	}

	runID := xid.New()
	ctx := context.Background()

	for i := 0; i < *ns; i++ {
		ns := fmt.Sprintf("%s-%d", runID, i)
		log.Printf("Creating namespace: %q", ns)
		_, err := clientset.CoreV1().Namespaces().Create(ctx, &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}, metav1.CreateOptions{})
		if err != nil {
			log.Fatalf("Error creating namespace: %v", err)
		}
		log.Printf("Created namespace: %q", ns)
	}
}
