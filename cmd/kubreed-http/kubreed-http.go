package main

import (
	"fmt"
	"kubreed/pkg/libs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	configJSON := os.Getenv(libs.ConfigEnvVar)
	c, err := libs.GetConfigFromJSON(configJSON)
	if err != nil {
		log.Fatalf("ENV variable not set properly for configuration: %v", err)
		return
	}

	log.Printf("Config is: %#v", c)

	// prepare server
	for i := 0; i < c.APICount; i++ {
		endpoint := fmt.Sprintf("/api%d", i)
		http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			randSleep := rand.Int63n((c.ResponseTime).Milliseconds())
			<-time.After(time.Duration(randSleep))
			w.Write([]byte("OK"))
			log.Printf("HTTPServer processed request from: %q", r.RemoteAddr)
		})
	}

	// launch client threads that talk to other servers
	go func() {
		reqCounter := 0

		for {
			// We can add a select loop here and gracefully exit if needed
			for _, svc := range c.RemoteServices {
				for apiIter := 0; apiIter < c.APICount; apiIter++ {
					go func(svc string, apiIter int) {
						url := fmt.Sprintf("http://%s/api%d", svc, apiIter)
						resp, err := http.Get(url)
						if err != nil {
							log.Printf("HTTPClient GET %q failed: %v", url, err)
						} else {
							log.Printf("HTTPClient GET %q: %v", url, resp.Status)
						}
						resp.Body.Close()
					}(svc, apiIter)
					reqCounter++

					if reqCounter == c.RPS {
						reqCounter = 0
						<-time.After(time.Second)
						log.Print("---------------------")
					}
				}
			}
		}
	}()

	// Launch server and block forever
	log.Fatal(http.ListenAndServe(":80", nil))
}
