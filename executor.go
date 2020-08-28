package main

import (
	log "./log"
	"os/exec"
	"syscall"
	//	"bytes"
	//	"strings"
	"bufio"
	"fmt"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"os/signal"
)

type Executor struct {
	path       string
	args       []string
	tty        *os.File
	rows, cols uint
	client     *Client
}

const DEBUG_EXEC = false
const DEBUG_EXEC_DUMP = false

func NewExecutor(client *Client, cols, rows uint, path string, args []string) *Executor {
	return &Executor{path: path, args: args, client: client, rows: rows, cols: cols}
}

func (executor *Executor) Execute() error {
	var err error

	_, err = pty.GetsizeFull(os.Stdin)

	cmd := exec.Command(executor.path, executor.args...)

	log.Debug("%s start", executor.Desc())

	oldSize, err := pty.GetsizeFull(os.Stdin)
	if err != nil {
		log.Error("%s fail pty getsize: %s", executor.Desc(), err)
	}

	var size = &pty.Winsize{Cols: uint16(executor.cols), Rows: uint16(executor.rows)}

	executor.tty, err = pty.StartWithSize(cmd, size)
	if err != nil {
		log.Error("%s fail pty start: %s", executor.Desc(), err)
		return log.NewError("fail pty start: %s", err)
	}
	defer func() {
		_ = executor.tty.Close()
	}()

	log.Debug("%s resize %dx%d", executor.Desc(), size.Cols, size.Rows)
	str := fmt.Sprintf("\033[8;%d;%dt", size.Rows, size.Cols)
	os.Stdout.Write([]byte(str))
	log.Debug("%s reset", executor.Desc())
	os.Stdout.Write([]byte("\033[H\033[2J"))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go executor.ProcessWindowChange(ch)
	ch <- syscall.SIGWINCH

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Error("%s error make raw: %s", executor.Desc(), err)
		return log.NewError("error make raw: %s", err)
	}
	defer func() {
		_ = terminal.Restore(int(os.Stdin.Fd()), oldState)
	}()

	go executor.CopyStdinToTTY()
	executor.ReadTTY()

	log.Debug("%s resize %dx%d", executor.Desc(), oldSize.Cols, oldSize.Rows)
	str = fmt.Sprintf("\033[8;%d;%dt", oldSize.Rows, oldSize.Cols)
	os.Stdout.Write([]byte(str))

	log.Debug("%s done", executor.Desc())
	return nil
}

func (executor *Executor) ReadTTY() {

	reader := bufio.NewReader(executor.tty)
	var buf []byte = make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Debug("%s read stdout error: %s", executor.Desc(), err)
			break
		}
		if DEBUG_EXEC_DUMP {
			log.Debug("%s read %d byte tty:\n%s", executor.Desc(), n, log.Dump(buf, n, 0))
		} else if DEBUG_EXEC {
			log.Debug("%s read %d byte tty", executor.Desc(), n)
		}
		os.Stdout.Write(buf[0:n])
		executor.client.SendText(buf[0:n])

	}
}

func (executor *Executor) ProcessWindowChange(ch chan os.Signal) {
	for range ch {
		str := ""
		rows, cols, err := pty.Getsize(os.Stdin)
		if err == nil {
			str = fmt.Sprintf("%dx%d", cols, rows)
		}
		log.Debug("window size %s", str)
	}
}

func (executor *Executor) CopyStdinToTTY() {
	var err error
	_, err = io.Copy(executor.tty, os.Stdin)
	if err != nil {
		log.Error("%s copy error: %s", executor.Desc(), err)
	}
}

func (executor *Executor) Desc() string {
	ret := "executor["
	ret += executor.path
	ret += "]"
	return ret
}
