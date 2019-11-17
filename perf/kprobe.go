package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	kp "github.com/zxdvd/go-libs/perf/kprobe"
)


type options struct {
	probe string
	uri     string
}

var cliOpt options

func init() {
	flag.StringVar(&cliOpt.probe, "probe", "", "probe string")
}

func doEvent(kp *kp.Kprobe, event string) {
	fmt.Println(event)
}

func kprobe(opt *options) error {
	quitCh := make(chan os.Signal, 1)
	stopCh := make(chan struct{}, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<- quitCh
		fmt.Printf("get signal, will quit soon\n")
		stopCh <- struct{}{}
	}()

	if err := kp.KprobeStream(opt.probe, doEvent, stopCh); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if cliOpt.probe == "" {
		fmt.Printf("please give a probe string like p:test1 sys_execve")
		os.Exit(-1)
	}
	if err := kprobe(&cliOpt); err != nil {
		panic(err)
	}
}
