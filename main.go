package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/valyala/fasthttp"
)

var start time.Time
var hostname string

func init() {
	hostname, _ = os.Hostname()
}

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "80"
	}

	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {

	switch string(ctx.Path()) {
	case "/echo":
		ctx.WriteString(hostname)
	default:
		w := ctx.Response.BodyWriter()

		cinfo, _ := cpu.Info()
		vm, _ := mem.VirtualMemory()

		host, _ := host.Info()

		fmt.Fprintln(w, "Environ \t "+strings.Join(os.Environ(), "-"))
		fmt.Fprintln(w, "Hostname \t "+hostname)

		fmt.Fprintln(w, "Headers")

		ctx.Request.Header.VisitAll(func(key, value []byte) {
			fmt.Fprintln(w, "\t \t", string(key), string(value))
		})

		fmt.Fprintln(w, "Request.Addr \t "+ctx.RemoteAddr().String())
		fmt.Fprintln(w, "RequestURI \t "+ctx.URI().String())

		fmt.Fprintln(w, "Uptime \t\t "+time.Now().Sub(start).String())

		fmt.Fprintf(w, "Memory \t\t total %s available %s used %s free %s percent %%%f \t\t \n",
			byteCountSI(vm.Total),
			byteCountSI(vm.Available),
			byteCountSI(vm.Used),
			byteCountSI(vm.Free),
			vm.UsedPercent)

		fmt.Fprintf(w, "Swap \t\t total %s free %s cached %s\n",
			byteCountSI(vm.SwapTotal),
			byteCountSI(vm.SwapFree),
			byteCountSI(vm.SwapCached))

		for _, cpu := range cinfo {
			fmt.Fprintf(w, "Cpu \t\t model %s core %d mhz %f\n",
				cpu.ModelName,
				cpu.Cores,
				cpu.Mhz)
		}

		fmt.Fprintf(w, "Host \t\t id %s name %s os %s vrole %s vsystem %s\n",
			host.HostID,
			host.Hostname,
			host.OS,
			host.VirtualizationRole,
			host.VirtualizationSystem)
	}
}

func byteCountSI(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
