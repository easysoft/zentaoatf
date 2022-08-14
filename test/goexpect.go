package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc/codes"
)

// DefaultTimeout is the default Expect timeout.
const DefaultTimeout = 60 * time.Second

const (
	checkDuration     = 2 * time.Second // checkDuration how often to check for new output.
	defaultBufferSize = 8192            // defaultBufferSize is the default io buffer size.
)
const (
	// OKTag marks the desired state was reached.
	OKTag = Tag(iota)
	// FailTag means reaching this state will fail the Switch/Case.
	FailTag
	// ContinueTag will recheck for matches.
	ContinueTag
	// NextTag skips match and continues to the next one.
	NextTag
	// NoTag signals no tag was set for this case.
	NoTag
)

// Status contains an errormessage and a status code.
type Status struct {
	code codes.Code
	msg  string
}

// GExpect implements the Expecter interface.
type GExpect struct {
	// cmd contains the cmd information for the spawned process.
	cmd *exec.Cmd
	// snd is the channel used by the Send command to send data into the spawned command.
	snd chan string
	// rcv is used to signal the Expect commands that new data arrived.
	rcv chan struct{}
	// chkMu lock protecting the check function.
	chkMu sync.RWMutex
	// chk contains the function to check if the spawned command is alive.
	chk func(*GExpect) bool
	// cls contains the function to close spawned command.
	cls func(*GExpect) error
	// timeout contains the default timeout for a spawned command.
	timeout time.Duration
	// sendTimeout contains the default timeout for a send command.
	sendTimeout time.Duration
	// chkDuration contains the duration between checks for new incoming data.
	chkDuration time.Duration
	// verbose enables verbose logging.
	verbose bool
	// verboseWriter if set specifies where to write verbose information.
	verboseWriter io.Writer
	// teeWriter receives a duplicate of the spawned process's output when set.
	teeWriter io.WriteCloser
	// PartialMatch enables the returning of unmatched buffer so that consecutive expect call works.
	partialMatch bool
	// bufferSize is the size of the io buffers in bytes.
	bufferSize int
	// bufferSizeIsSet tracks whether the bufferSize was set for a given GExpect instance.
	bufferSizeIsSet bool

	// mu protects the output buffer. It must be held for any operations on out.
	mu  sync.Mutex
	out bytes.Buffer
}

var stdout bytes.Buffer

func Spawn(command string, timeout time.Duration) (*GExpect, <-chan error, error) {
	return SpawnWithArgs(strings.Fields(command), timeout)
}
func SpawnWithArgs(command []string, timeout time.Duration) (*GExpect, <-chan error, error) {

	if timeout < 1 {
		timeout = DefaultTimeout
	}
	// Get the command up and running
	cmd := exec.Command(command[0], command[1:]...)
	// This ties the commands Stdin,Stdout & Stderr to the virtual terminal we created
	cmd.Stdin, cmd.Stdout, cmd.Stderr = &stdout, &stdout, &stdout
	// New process needs to be the process leader and control of a tty
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true}
	e := &GExpect{
		rcv:         make(chan struct{}),
		snd:         make(chan string),
		cmd:         cmd,
		timeout:     timeout,
		chkDuration: checkDuration,
		cls: func(e *GExpect) error {
			if e.cmd != nil {
				return e.cmd.Process.Kill()
			}
			return nil
		},
		chk: func(e *GExpect) bool {
			if e.cmd.Process == nil {
				return false
			}
			// Sending Signal 0 to a process returns nil if process can take a signal , something else if not.
			return e.cmd.Process.Signal(syscall.Signal(0)) == nil
		},
	}

	// Set the buffer size to the default if expect.BufferSize(...) is not utilized.
	if !e.bufferSizeIsSet {
		e.bufferSize = defaultBufferSize
	}

	res := make(chan error, 1)
	go e.runcmd(res)
	// Wait until command started
	return e, res, <-res
}

// runcmd executes the command and Wait for the return value.
func (e *GExpect) runcmd(res chan error) {
	if err := e.cmd.Start(); err != nil {
		res <- err
		return
	}
	// Moving the go read/write functions here makes sure the command is started before first checking if it's running.
	clean := make(chan struct{})
	chDone := e.goIO(clean)
	// Signal command started
	res <- nil
	cErr := e.cmd.Wait()
	close(chDone)
	stdout.Reset()
	// make sure the read/send routines are done before closing the pty.
	<-clean
	res <- cErr
}

// goIO starts the io handlers.
func (e *GExpect) goIO(clean chan struct{}) (done chan struct{}) {
	done = make(chan struct{})
	var ptySync sync.WaitGroup
	ptySync.Add(2)
	go e.read(done, &ptySync)
	go e.send(done, &ptySync)
	go func() {
		ptySync.Wait()
		stdout.Reset()
		close(clean)
	}()
	return done
}

// read reads from the PTY master and forwards to active Expect function.
func (e *GExpect) read(done chan struct{}, ptySync *sync.WaitGroup) {
	defer ptySync.Done()
	buf := make([]byte, e.bufferSize)
	for {
		nr, err := stdout.Read(buf)
		if err != nil && !e.check() {
			if e.teeWriter != nil {
				e.teeWriter.Close()
			}
			if err == io.EOF {
				if e.verbose {
					log.Printf("read closing down: %v", err)
				}
				return
			}
			return
		}
		// Tee output to writer
		if e.teeWriter != nil {
			e.teeWriter.Write(buf[:nr])
		}
		// Add to buffer
		e.mu.Lock()
		e.out.Write(buf[:nr])
		e.mu.Unlock()
		// Ping Expect function
		select {
		case e.rcv <- struct{}{}:
		default:
		}
	}
}
func (e *GExpect) check() bool {
	e.chkMu.RLock()
	defer e.chkMu.RUnlock()
	return e.chk(e)
}

// Send sends a string to spawned process.
func (e *GExpect) Send(in string) error {
	if !e.check() {
		return errors.New("expect: Process not running")
	}

	if e.sendTimeout == 0 {
		e.snd <- in
	} else {
		select {
		case <-time.After(e.sendTimeout):
			return fmt.Errorf("send to spawned process command reached the timeout %v", e.sendTimeout)
		case e.snd <- in:
		}
	}

	if e.verbose {
		if e.verboseWriter != nil {
			vStr := fmt.Sprintln("Sent:" + fmt.Sprintf(" %q", in))
			_, err := e.verboseWriter.Write([]byte(vStr))
			if err != nil {
				log.Printf("Write to Verbose Writer failed: %v", err)
			}
		} else {
			log.Printf("Sent: %q", in)
		}
	}

	return nil
}

// send writes to the PTY master.
func (e *GExpect) send(done chan struct{}, ptySync *sync.WaitGroup) {
	defer ptySync.Done()
	for {
		select {
		case <-done:
			return
		case sstr, ok := <-e.snd:
			if !ok {
				return
			}
			if _, err := stdout.Write([]byte(sstr)); err != nil || !e.check() {
				log.Printf("send failed: %v", err)
			}
		}
	}
}

// Read implements the reader interface for the out buffer.
func (e *GExpect) Read(p []byte) (nr int, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.out.Read(p)
}

// Close closes the Spawned session.
func (e *GExpect) Close() error {
	return e.cls(e)
}

type Tag int32

// NewStatus creates a Status with the provided code and message.
func NewStatus(code codes.Code, msg string) *Status {
	return &Status{code, msg}
}

// Case used by the ExpectSwitchCase to take different Cases.
// Implements the Caser interface.
type Case struct {
	// R is the compiled regexp to match.
	R *regexp.Regexp
	// S is the string to send if Regexp matches.
	S string
	// T is the Tag for this Case.
	T func() (Tag, *Status)
	// Rt specifies number of times to retry, only used for cases tagged with Continue.
	Rt int
}

// Tag returns the tag for this case.
func (c *Case) Tag() (Tag, *Status) {
	if c.T == nil {
		return NoTag, NewStatus(codes.OK, "no Tag set")
	}
	return c.T()
}

// RE returns the compiled regular expression.
func (c *Case) RE() (*regexp.Regexp, error) {
	return c.R, nil
}

// Caser is an interface for ExpectSwitchCase and Batch to be able to handle
// both the Case struct and the more script friendly BCase struct.
type Caser interface {
	// RE returns a compiled regexp
	RE() (*regexp.Regexp, error)
	// Send returns the send string
	String() string
	// Tag returns the Tag.
	Tag() (Tag, *Status)
	// Retry returns true if there are retries left.
	Retry() bool
}

// Retry decrements the Retry counter and checks if there are any retries left.
func (c *Case) Retry() bool {
	defer func() { c.Rt-- }()
	return c.Rt > 0
}

// Send returns the string to send if regexp matches
func (c *Case) String() string {
	return c.S
}

// Expect reads spawned processes output looking for input regular expression.
// Timeout set to 0 makes Expect return the current buffer.
// Negative timeout value sets it to Default timeout.
func (e *GExpect) Expect(re *regexp.Regexp, timeout time.Duration) (string, []string, error) {
	out, match, _, err := e.ExpectSwitchCase([]Caser{&Case{re, "", nil, 0}}, timeout)
	return out, match, err
}
func (e *GExpect) returnUnmatchedSuffix(p string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	newBuffer := bytes.NewBufferString(p)
	newBuffer.WriteString(e.out.String())
	e.out = *newBuffer
}

// Err is a helper to handle errors.
func (s *Status) Err() error {
	if s == nil || s.code == codes.OK {
		return nil
	}
	return s
}

// Error is here to adhere to the error interface.
func (s *Status) Error() string {
	return s.msg
}

// NewStatusf returns a Status with the provided code and a formatted message.
func NewStatusf(code codes.Code, format string, a ...interface{}) *Status {
	return NewStatus(code, fmt.Sprintf(fmt.Sprintf(format, a...)))
}

// TimeoutError is the error returned by all Expect functions upon timer expiry.
type TimeoutError int

// Error implements the Error interface.
func (t TimeoutError) Error() string {
	return fmt.Sprintf("expect: timer expired after %d seconds", time.Duration(t)/time.Second)
}

// ExpectSwitchCase checks each Case against the accumulated out buffer, sending specified
// string back. Leaving Send empty will Send nothing to the process.
// Substring expansion can be used eg.
// 	Case{`vf[0-9]{2}.[a-z]{3}[0-9]{2}\.net).*UP`,`show arp \1`}
// 	Given: vf11.hnd01.net            UP      35 (4)        34 (4)          CONNECTED         0              0/0
// 	Would send: show arp vf11.hnd01.net
func (e *GExpect) ExpectSwitchCase(cs []Caser, timeout time.Duration) (string, []string, int, error) {
	// Compile all regexps
	rs := make([]*regexp.Regexp, 0, len(cs))
	for _, c := range cs {
		re, err := c.RE()
		if err != nil {
			return "", []string{""}, -1, err
		}
		rs = append(rs, re)
	}
	// Setup timeouts
	// timeout == 0 => Just dump the buffer and exit.
	// timeout < 0  => Set default value.
	if timeout < 0 {
		timeout = e.timeout
	}
	timer := time.NewTimer(timeout)
	check := e.chkDuration
	// Check if any new data arrived every checkDuration interval.
	// If timeout/4 is less than the checkout interval we set the checkout to
	// timeout/4. If timeout ends up being 0 we bump it to one to keep the Ticker from
	// panicking.
	// All this b/c of the unreliable channel send setup in the read function,making it
	// possible for Expect* functions to miss the rcv signal.
	//
	// from read():
	//		// Ping Expect function
	//		select {
	//		case e.rcv <- struct{}{}:
	//		default:
	//		}
	//
	// A signal is only sent if any Expect function is running. Expect could miss it
	// while playing around with buffers and matching regular expressions.
	if timeout>>2 < check {
		check = timeout >> 2
		if check <= 0 {
			check = 1
		}
	}
	chTicker := time.NewTicker(check)
	defer chTicker.Stop()
	// Read in current data and start actively check for matches.
	var tbuf bytes.Buffer
	if _, err := io.Copy(&tbuf, e); err != nil {
		return tbuf.String(), nil, -1, fmt.Errorf("io.Copy failed: %v", err)
	}
	for {
	L1:
		for i, c := range cs {
			if rs[i] == nil {
				continue
			}
			match := rs[i].FindStringSubmatch(tbuf.String())
			if match == nil {
				continue
			}

			t, s := c.Tag()
			if t == NextTag && !c.Retry() {
				continue
			}

			if e.verbose {
				if e.verboseWriter != nil {
					vStr := fmt.Sprintln("Match for RE:" + fmt.Sprintf(" %q found: %q Buffer: %s", rs[i].String(), match, tbuf.String()))
					for n, bytesRead, err := 0, 0, error(nil); bytesRead < len(vStr); bytesRead += n {
						n, err = e.verboseWriter.Write([]byte(vStr)[n:])
						if err != nil {
							log.Printf("Write to Verbose Writer failed: %v", err)
							break
						}
					}
				} else {
					log.Printf("Match for RE: %q found: %q Buffer: %q", rs[i].String(), match, tbuf.String())
				}
			}

			tbufString := tbuf.String()
			o := tbufString

			if e.partialMatch {
				// Return the part of the buffer that is not matched by the regular expression so that the next expect call will be able to match it.
				matchIndex := rs[i].FindStringIndex(tbufString)
				o = tbufString[0:matchIndex[1]]
				e.returnUnmatchedSuffix(tbufString[matchIndex[1]:])
			}

			tbuf.Reset()

			st := c.String()
			// Replace the submatches \[0-9]+ in the send string.
			if len(match) > 1 && len(st) > 0 {
				for i := 1; i < len(match); i++ {
					// \(submatch) will be expanded in the Send string.
					// To escape use \\(number).
					si := strconv.Itoa(i)
					r := strings.NewReplacer(`\\`+si, `\`+si, `\`+si, `\\`+si)
					st = r.Replace(st)
					st = strings.Replace(st, `\\`+si, match[i], -1)
				}
			}
			// Don't send anything if string is empty.
			if st != "" {
				if err := e.Send(st); err != nil {
					return o, match, i, fmt.Errorf("failed to send: %q err: %v", st, err)
				}
			}
			// Tag handling.
			switch t {
			case OKTag, FailTag, NoTag:
				return o, match, i, s.Err()
			case ContinueTag:
				if !c.Retry() {
					return o, match, i, s.Err()
				}
				break L1
			case NextTag:
				break L1
			default:
				s = NewStatusf(codes.Unknown, "Tag: %d unknown, err: %v", t, s)
			}
			return o, match, i, s.Err()
		}
		if !e.check() {
			nr, err := io.Copy(&tbuf, e)
			if err != nil {
				return tbuf.String(), nil, -1, fmt.Errorf("io.Copy failed: %v", err)
			}
			if nr == 0 {
				return tbuf.String(), nil, -1, errors.New("expect: Process not running")
			}
		}
		select {
		case <-timer.C:
			// Expect timeout.
			nr, err := io.Copy(&tbuf, e)
			if err != nil {
				return tbuf.String(), nil, -1, fmt.Errorf("io.Copy failed: %v", err)
			}
			// If we got no new data we return otherwise give it another chance to match.
			if nr == 0 {
				return tbuf.String(), nil, -1, TimeoutError(timeout)
			}
			timer = time.NewTimer(timeout)
		case <-chTicker.C:
			// Periodical timer to make sure data is handled in case the <-e.rcv channel
			// was missed.
			if _, err := io.Copy(&tbuf, e); err != nil {
				return tbuf.String(), nil, -1, fmt.Errorf("io.Copy failed: %v", err)
			}
		case <-e.rcv:
			// Data to fetch.
			nr, err := io.Copy(&tbuf, e)
			if err != nil {
				return tbuf.String(), nil, -1, fmt.Errorf("io.Copy failed: %v", err)
			}
			// timer shoud be reset when new output is available.
			if nr > 0 {
				timer = time.NewTimer(timeout)
			}
		}
	}
}
