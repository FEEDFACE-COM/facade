package main

import (
	facade "./facade"
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
	"path"
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

	log.Info("%s start", executor.Desc())

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

	if DEBUG_EXEC {
		log.Debug("%s resize %dx%d", executor.Desc(), size.Cols, size.Rows)
	}
	str := fmt.Sprintf("\033[8;%d;%dt", size.Rows, size.Cols)
	os.Stdout.Write([]byte(str))
	if DEBUG_EXEC {
		log.Debug("%s reset", executor.Desc())
	}
	os.Stdout.Write([]byte("\033[H\033[2J"))

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go executor.ProcessWindowChange(ch)
	//	ch <- syscall.SIGWINCH

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Error("%s error make raw: %s", executor.Desc(), err)
		return log.NewError("error make raw: %s", err)
	}
	defer func() {
		_ = terminal.Restore(int(os.Stdin.Fd()), oldState)
	}()

	executor.client.SendText([]byte("\033[H\033[2J"))

	go executor.CopyStdinToTTY()
	executor.ReadTTY()

	if DEBUG_EXEC {
		log.Debug("%s resize %dx%d", executor.Desc(), oldSize.Cols, oldSize.Rows)
	}
	str = fmt.Sprintf("\033[8;%d;%dt", oldSize.Rows, oldSize.Cols)
	os.Stdout.Write([]byte(str))

	if DEBUG_EXEC {
		log.Debug("%s done", executor.Desc())
	}
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
			log.Error("%s read stdout error: %s", executor.Desc(), err)
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
		var err error
		rows, cols, err := pty.Getsize(os.Stdin)
		if err != nil {
			log.Error("%s fail get size: %s", executor.Desc(), err)
			continue
		}
		err = pty.Setsize(executor.tty, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
		if err != nil {
			log.Error("%s fail inherit size: %s", executor.Desc(), err)
			continue
		}
		log.Info("%s window resized %dx%d", executor.Desc(), cols, rows)
		grid := facade.GridConfig{Width: uint64(cols), SetWidth: true, Height: uint64(rows), SetHeight: true}
		conf := facade.Config{Terminal: &facade.TermConfig{Grid: &grid}}
		executor.client.SendConf(&conf)
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
	ret += path.Base(executor.path)
	ret += "]"
	return ret
}
