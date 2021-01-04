package main

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"io"
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

	printInfo(os.Stdout)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		printInfo(w)
		fmt.Fprint(w, "\n")
		fmt.Fprintln(w, "Request.Host \t "+r.Host)
		fmt.Fprintln(w, "Request.Addr \t "+r.RemoteAddr)

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

func printInfo(w io.Writer) {
	v, _ := mem.VirtualMemory()
	hostinfo, _ := host.Info()

	minfo := fmt.Sprintf("Total: %v, Free:%v, UsedPercent:%f%%", ByteCountDecimal(v.Total), ByteCountDecimal(v.Free), v.UsedPercent)

	info := map[string]string{
		"Environ":       strings.Join(os.Environ(), " - "),
		"Go.Version":    runtime.Version(),
		"Go.NumCpu":     strconv.Itoa(runtime.NumCPU()),
		"Host.Id":       hostinfo.HostID,
		"Host.Os":       hostinfo.OS,
		"Host.Hostname": hostinfo.Hostname,
		"Host.Uptime":   strconv.FormatUint(hostinfo.Uptime, 10),
		"Host.Procs":    strconv.FormatUint(hostinfo.Procs, 10),
		"Host.Platform": fmt.Sprintf("%s %s", hostinfo.Platform, hostinfo.PlatformVersion),
		"Memory":        minfo,
	}

	var keys []string
	for k := range info {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		_, _ = fmt.Fprintln(w, k, " \t ", info[k])
	}
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
