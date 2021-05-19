package main

import (
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

func main() {
	port := 80

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

func ByteCountDecimal(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
