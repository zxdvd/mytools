package main

import (
	"regexp"
	"strings"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/zxdvd/go-libs/perf/kprobe"
)

// copy from brendangregg's execsnoop
const kprobeCmd = "p:execsnoop sys_execve arg1=+0(+0(%si)):string arg2=+0(+8(%si)):string arg3=+0(+16(%si)):string arg4=+0(+24(%si)):string arg5=+0(+32(%si)):string arg6=+0(+40(%si)):string arg7=+0(+48(%si)):string arg8=+0(+56(%si)):string arg9=+0(+64(%si)):string"


var pat = regexp.MustCompile(" arg[0-9]=")

func doEvent(kp *kprobe.Kprobe, event string) {
	parts := pat.Split(event, -1)
	for i, part := range parts {
		if i == 0 {
			if !strings.Contains(part, kp.Name()) {
				return
			}
			part = strings.TrimLeft(part, " ")
			subparts := strings.Split(part, " ")
			if len(subparts) > 0 {
				fmt.Printf("%s  ", subparts[0])
			}
			continue
		}
		if strings.HasPrefix(part, "(fault)") {
			break
		}
		fmt.Printf("%s ", part[1:len(part)-1])
	}
	fmt.Printf("\n")
}

func execsnoop() error {
	quitCh := make(chan os.Signal, 1)
	stopCh := make(chan struct{}, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<- quitCh
		fmt.Printf("get signal, will quit soon\n")
		stopCh <- struct{}{}
	}()

	if err := kprobe.KprobeStream(kprobeCmd, doEvent, stopCh); err != nil {
		return err
	}
	return nil
}

func main() {
	execsnoop()
}
