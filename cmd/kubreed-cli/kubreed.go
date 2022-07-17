package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"kubreed/pkg/libs"

	"github.com/rs/xid"
	flag "github.com/spf13/pflag"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// default values
const (
	Namespaces  = 1
	Deployments = 5
	Pods        = 3
	APIs        = 10
	RPS         = 1
	Branching   = 4
	Latency     = time.Second * 2
)

func main() {
	ns := flag.IntP("namespaces", "n", Namespaces, "Number of Namespaces to create")
	deps := flag.IntP("deployments", "d", Deployments, "Number of Deployments/Services to create per Namespace")
	pods := flag.Int32P("pods", "p", Pods, "Number of Pods to create per Deployment")
	api := flag.IntP("apis", "a", APIs, "Number of APIs per Pod")
	rps := flag.IntP("rps", "r", RPS, "Outgoing rps by each client Pod")
	branching := flag.IntP("branching", "b", Branching, "Number of Services to which each client Pod should make requests")
	latency := flag.DurationP("latency", "l", Latency, "Maximum response time in milliseconds for each API call")

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
		_, err = clientset.CoreV1().Namespaces().Create(ctx,
			&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}},
			metav1.CreateOptions{})
		if err != nil {
			log.Fatalf("Error creating namespace: %v", err)
			return
		}
		log.Printf("Created namespace: %q", ns)

		for j := 0; j < *deps; j++ {
			svcName := fmt.Sprintf("svc-%d", j)
			dep := fmt.Sprintf("dep-%d", j)
			labels := map[string]string{
				"app": dep,
			}
			objectMeta := metav1.ObjectMeta{
				Name:      dep,
				Namespace: ns,
				Labels:    labels,
			}

			log.Printf("Creating Deployment: %q", dep)
			_, err = clientset.AppsV1().Deployments(ns).Create(ctx,
				&appsv1.Deployment{
					ObjectMeta: objectMeta,
					Spec: appsv1.DeploymentSpec{
						Replicas: pods,
						Selector: &metav1.LabelSelector{
							MatchLabels: labels,
						},
						Template: v1.PodTemplateSpec{
							ObjectMeta: objectMeta,
							Spec: v1.PodSpec{
								Containers: []v1.Container{{
									Name:  "kubreed-http",
									Image: "psankar/kubreed-http:f25a398",
									Ports: []v1.ContainerPort{{
										ContainerPort: 80,
										Protocol:      "TCP",
									}},
									Env: []v1.EnvVar{{
										Name: libs.ConfigEnvVar,
										Value: `{
											"apiCount": 3,
											"responseTime": "1s",
											"rps": 10,
											"remoteServices": [
											  "svc1",
											  "svc2"
											]
										  }`,
									}},
								}},
							},
						},
					},
				},
				metav1.CreateOptions{})
			if err != nil {
				log.Fatalf("Error creating deployment: %v", err)
				return
			}
			log.Printf("Created deployment: %q", dep)

			log.Printf("Creating service: %q", svcName)
			_, err = clientset.CoreV1().Services(ns).Create(ctx,
				&v1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      svcName,
						Namespace: ns,
						Labels:    map[string]string{},
					},
					Spec: v1.ServiceSpec{
						Selector: map[string]string{
							"app": dep,
						},
						Ports: []v1.ServicePort{
							{
								Port: 80,
								TargetPort: intstr.IntOrString{
									Type:   intstr.Int,
									IntVal: 80,
								},
							},
						},
					},
				},
				metav1.CreateOptions{})
			if err != nil {
				log.Fatalf("Error creating service: %v", err)
			}
		}
	}
}
