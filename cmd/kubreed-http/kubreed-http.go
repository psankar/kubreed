package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/psankar/kubreed/pkg/libs"
	flag "github.com/spf13/pflag"
)

const (
	minDuration = time.Millisecond * 10
	maxDuration = time.Minute * 2
	minAPI      = 1
	maxAPI      = 20
)

func main() {
	podName := flag.String("podName", "", "Name of the Pod")
	apiCount := flag.Int("apiCount", libs.APIs, "Number of APIs")
	respTime := flag.Duration("respTime", libs.Latency, "Maximum response time in milliseconds for each API call")
	rps := flag.IntP("rps", "r", libs.RPS, "Outgoing rps from each client Pod")
	remoteServices := flag.StringArray("remoteServices", []string{}, "Remote services to talk to")

	flag.Parse()

	// validate configuration
	if len(*podName) == 0 {
		log.Fatal("Podname not passed")
		return
	}

	if *apiCount < minAPI || *apiCount > maxAPI {
		log.Fatal("API Count should be more than 1")
		return
	}

	if *respTime < minDuration || *respTime > maxDuration {
		log.Fatalf("respTime should be between %v and %v", minDuration, maxDuration)
		return
	}

	// prepare server
	for i := 0; i < *apiCount; i++ {
		endpoint := fmt.Sprintf("/api%d", i)
		http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			randSleep := rand.Int63n((*respTime).Milliseconds())
			<-time.After(time.Duration(randSleep))
			w.Write([]byte("OK"))
		})
	}

	// launch client threads that talk to other servers
	go func() {
		reqCounter := 0

		for {
			// We can add a select loop here and gracefully exit if needed
			for _, svc := range *remoteServices {
				for apiIter := 0; apiIter < *apiCount; apiIter++ {
					go func(svc string, apiIter int) {
						log.Printf("Make the call to http://%s/api%d", svc, apiIter)
					}(svc, apiIter)
					reqCounter++

					if reqCounter == *rps {
						reqCounter = 0
						<-time.After(time.Second)
						log.Print("---------------------")
					}
				}
			}
		}
	}()

	// Launch server and block forever
	log.Fatal(http.ListenAndServe(":8080", nil))
}
