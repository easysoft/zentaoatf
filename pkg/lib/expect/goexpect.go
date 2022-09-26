package expectUtils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// DefaultTimeout is the default Expect timeout.
const DefaultTimeout = 60 * time.Second

// GExpect implements the Expecter interface.
type GExpect struct {
	// cmd contains the cmd information for the spawned process.
	cmd *exec.Cmd
	out io.ReadCloser
	in  io.WriteCloser
	err io.ReadCloser
}

func Spawn(cmdStr string, timeout time.Duration) (expect *GExpect, err error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmdStr = strings.ReplaceAll(cmdStr, "/", "\\")
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}
	fmt.Println(cmd.String())

	if cmd == nil {
		err = errors.New("cmd is nil")
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	expect = &GExpect{
		cmd: cmd,
		out: stdout,
		in:  stdin,
		err: stderr,
	}
	return
}
func (e *GExpect) Expect(expect *regexp.Regexp, timeout time.Duration) (out string, err error) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	c := make(chan int, 1)
	go e.expectActual(c, expect, &out, &err)
	for {
		select {
		case <-c:
			return out, err
		case <-timer.C:
			if out != "" {
				err = errors.New(out)
			} else {
				err = errors.New("Time out")
			}
			return
		}
	}
}
func (e *GExpect) expectActual(c chan int, expect *regexp.Regexp, out *string, err *error) {
	reader1 := bufio.NewReader(e.out)
	for {
		line, err2 := reader1.ReadString('\n')
		if err2 != nil {
			err = &err2
			return
		}
		*out = fmt.Sprintf("%s%s", *out, line)
		if expect.MatchString(*out) {
			c <- 1
			return
		}
		if *err != nil || io.EOF == *err {
			break
		}
	}
}
func (e *GExpect) Send(msg string) (err error) {
	e.in.Write([]byte(msg))
	return
}
func (e *GExpect) Close() (err error) {
	e.cmd.Process.Kill()
	return
}
