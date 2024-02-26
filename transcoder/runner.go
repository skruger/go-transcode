package transcoder

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

const (
	RunnerStdErr = iota
	RunnerStdOut
	RunnerMetadata
)

type RunnerOutput struct {
	Data   string
	Source int
}

type TranscodeRunner struct {
	cmd     *exec.Cmd
	outputs chan RunnerOutput
	outDone chan interface{}
	errDone chan interface{}
	done    bool
}

func NewTranscodeRunner(cmd *exec.Cmd) *TranscodeRunner {
	return &TranscodeRunner{
		cmd:     cmd,
		outputs: make(chan RunnerOutput, 1000),
		outDone: make(chan interface{}),
		errDone: make(chan interface{}),
		done:    false,
	}
}

func (t *TranscodeRunner) Start() error {
	var outPipe io.ReadCloser
	var errPipe io.ReadCloser
	var err error
	outPipe, err = t.cmd.StdoutPipe()
	if err != nil {
		t.outputs <- RunnerOutput{
			Data:   fmt.Sprintf("error opening stdout pipe: %s\n", err),
			Source: RunnerMetadata,
		}
		return err
	}
	errPipe, err = t.cmd.StderrPipe()
	if err != nil {
		t.outputs <- RunnerOutput{
			Data:   fmt.Sprintf("error opening stderr pipe: %s\n", err),
			Source: RunnerMetadata,
		}
		return err
	}
	go func() {
		t.reader(outPipe, RunnerStdOut, "stdout")
		close(t.outDone)
	}()
	go func() {
		t.reader(errPipe, RunnerStdErr, "stderr")
		close(t.errDone)
	}()

	err = t.cmd.Start()
	if err != nil {
		t.outputs <- RunnerOutput{
			Data:   fmt.Sprintf("Run error: %s", err),
			Source: RunnerMetadata,
		}
		return err
	}

	return nil
}

func (t *TranscodeRunner) reader(pipe io.ReadCloser, source int, sourceName string) {
	for {
		// Wait for process to start
		if t.cmd.Process != nil {
			break
		}
	}
	scanner := bufio.NewScanner(pipe)
	scanner.Split(progressScanner)
	for scanner.Scan() {
		line := scanner.Text()
		t.outputs <- RunnerOutput{
			Data:   line,
			Source: source,
		}
	}

}

func (t *TranscodeRunner) ReceiveLine() (*RunnerOutput, bool) {
	for {
		if len(t.outputs) == 0 && t.done {
			return nil, true
		}
		if len(t.outputs) > 0 {
			output := <-t.outputs
			return &output, false
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (t *TranscodeRunner) Wait() error {
	err := t.cmd.Wait()
	if err != nil {
		t.outputs <- RunnerOutput{
			Data:   fmt.Sprintf("wait error: %s", err),
			Source: RunnerMetadata,
		}
	}
	<-t.errDone
	<-t.outDone
	t.done = true
	return err
}

func progressScanner(data []byte, atEOF bool) (advance int, token []byte, err error) {
	r := bytes.IndexByte(data, '\r')
	n := bytes.IndexByte(data, '\n')
	if r == -1 && n == -1 {
		if !atEOF {
			return 0, nil, nil
		}
		return 0, data, bufio.ErrFinalToken
	}
	if r+1 != n {
		return r + 1, data[:r+1], nil
	}
	return n + 1, data[:n+1], nil
}

func (r *RunnerOutput) IsStatus() bool {
	return len(r.Data) > 0 && r.Data[len(r.Data)-1] == '\r'
}
