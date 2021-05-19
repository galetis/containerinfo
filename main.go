package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var port int

func init()  {
	flag.IntVar(&port, "port", lookupEnvOrInt("PORT", 80), "port")
	flag.Parse()
}

func main() {

	i, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		port = i
	}

	start := time.Now()

	log.Printf("Listening on: %d", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintln(w, "Environ \t "+strings.Join(os.Environ(), "-"))

		fmt.Fprintln(w, "Headers")

		keys := []string{}
		for k := range r.Header {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintln(w, "\t \t", k, r.Header.Get(k))
		}

		fmt.Fprintln(w, "Request.Addr \t "+r.RemoteAddr)
		fmt.Fprintln(w, "RequestURI \t "+r.RequestURI)

		fmt.Fprintln(w, "Uptime \t\t "+time.Now().Sub(start).String())
	})

	http.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan int)
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for {
					select {
					case <-done:
						return
					default:
					}
				}
			}()
		}
		time.Sleep(time.Second * 5)
		close(done)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func lookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}