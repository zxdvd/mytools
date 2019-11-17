package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/zxdvd/go-libs/perf/uprobe"
)


type options struct {
	probe string
	uri     string
}

var cliOpt options

func init() {
	flag.StringVar(&cliOpt.probe, "probe", "", "probe string")
}

func doEvent(up *uprobe.Uprobe, event string) {
	fmt.Println(event)
}

func run(opt *options) error {
	quitCh := make(chan os.Signal, 1)
	stopCh := make(chan struct{}, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<- quitCh
		fmt.Printf("get signal, will quit soon\n")
		stopCh <- struct{}{}
	}()

	if err := uprobe.UprobeStream(opt.probe, doEvent, stopCh); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if cliOpt.probe == "" {
		fmt.Printf("please give a probe string like p:test1 /bin/bash:readline\n")
		os.Exit(-1)
	}
	if err := run(&cliOpt); err != nil {
		panic(err)
	}
}
