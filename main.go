package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/h4yfans/adjust/tool"
)

const DefaultConcurrentLimit = 10

var (
	parallel = flag.Uint("parallel", DefaultConcurrentLimit, "number of parallel processes")
	timeout  = flag.Duration("timeout", 30*time.Second, "timeout duration for http request")
)

var usage = `Usage: ./myhttp [options...] [urls...]

Options:
	-parallel Number of workers to run concurrently. Default is 10.
	-timeout  Timeout for each request in seconds. Default is 30.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}
	flag.Parse()

	httpClient := &http.Client{
		Timeout: *timeout,
	}
	if err := tool.NewTool(*parallel, httpClient).Run(flag.Args()...); err != nil {
		usageAndExit(err.Error())
	}
}

func usageAndExit(msg string) {

	fmt.Fprintf(os.Stderr, msg)
	fmt.Fprintf(os.Stderr, "\n\n")
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
