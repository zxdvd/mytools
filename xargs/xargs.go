package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/zxdvd/go-libs/std-helper/shlex"
)

var verbose bool

func init () {
	flag.BoolVar(&verbose, "verbose", false, "verbose, show command")
}

func main() {
	flag.Parse()
	args := flag.Args()
	cmdFmt := strings.Join(args, " ")
	fmt.Printf("cmd is %s\n", cmdFmt)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF && len(line) == 0{
			break
		}
		line = line[:len(line)-1]
		cmd := fmt.Sprintf(cmdFmt, shlex.Quote(string(line)))
		if verbose {
			fmt.Printf("CMD: %s\n", cmd)
		}
		out, _ := exec.Command("sh", "-c", cmd).Output()
		fmt.Printf(string(out))
	}
}
