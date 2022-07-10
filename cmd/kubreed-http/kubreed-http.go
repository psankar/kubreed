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
	respTime := flag.Duration("respTime", libs.RespTime, "Maximum response time in milliseconds for each API call")
	rps := flag.IntP("rps", "r", libs.RPS, "Number of Requests Per Second made by each client Pod")
	branching := flag.IntP("branching", "b", libs.Branching, "Number of Services to which each client Pod should make requests")

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
		sleepTimer := time.Second * 5
		timer := time.NewTimer(sleepTimer)
		for {
			go func(reqCounter int) {
				for j := 0; j < *branching; j++ {
					log.Printf("reqCounter %d GET CALL: %d", reqCounter, j)
				}
			}(reqCounter)
			reqCounter++

			if reqCounter == *rps {
				<-timer.C
				reqCounter = 0
				timer.Reset(sleepTimer)
				log.Printf("New set of client requests started after time")
			}
		}
	}()

	// Launch server and block forever
	log.Fatal(http.ListenAndServe(":8080", nil))
}
