package expectUtils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"runtime"
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
	// var stdout, stdin, stderr bytes.Buffer
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	if cmd == nil {
		err = errors.New("cmd is nil")
		return
	}
	// cmd.Stdin, cmd.Stdout, cmd.Stderr = &stdout, &stdin, &stderr
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
		fmt.Println(00, err)
	}

	expect = &GExpect{
		cmd: cmd,
		out: stdout,
		in:  stdin,
		err: stderr,
	}

	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Println(01, err)
	// }
	return
}
func (e *GExpect) Expect(expect *regexp.Regexp, timeout time.Duration) (out string, match string, err error) {
	reader1 := bufio.NewReader(e.out)
	for true {
		line, err := reader1.ReadString('\n')
		out = fmt.Sprintf("%s%s", out, line)
		if expect.MatchString(out) {
			return out, match, err
		}
		if err != nil || io.EOF == err {
			break
		}
	}
	return
}
func (e *GExpect) Send(msg string) (err error) {
	e.in.Write([]byte(msg))
	return
}
func (e *GExpect) Close() (err error) {
	e.cmd.Process.Kill()
	return
}
