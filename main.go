package main

import (
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/http"
)

func main() {
	var port int

	flag.IntVar(&port, "port", 8080, "listen port")

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		v, _ := mem.VirtualMemory()
		host, _ := host.Info()

		_, _ = fmt.Fprintf(w, "Memory: Total: %v, Free:%v, UsedPercent:%f%%\n", ByteCountDecimal(v.Total), ByteCountDecimal(v.Free), v.UsedPercent)

		_, _ = fmt.Fprintln(w, "Host:", host)

		_, _ = fmt.Fprintln(w, "Host:", r.Host)
		_, _ = fmt.Fprintln(w, "RemoteAddr:", r.RemoteAddr)

		_, _ = fmt.Fprintln(w, r.Method, r.RequestURI)
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
