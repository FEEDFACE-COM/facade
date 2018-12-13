package main

import(
	log "./log"
	"os/exec"
	"syscall"
//	"bytes"
//	"strings"
	"os"
	"io"
	"bufio"
	"fmt"
	"os/signal"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"

)

type Executor struct{
	path string
	args []string
}


func NewExecutor(path string, args []string) *Executor { 
	return &Executor{path:path, args:args} 
}


func (executor *Executor) Execute() error {
	
	

	cmd := exec.Command(executor.path, executor.args ...)
	
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Error("fail pty start: %s",err)
		return log.NewError("fail pty start: %s",err)
	}
	
	
	defer func() { _ = ptmx.Close() }() 
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			err := pty.InheritSize(os.Stdin,ptmx); 
			if err != nil {
				log.Error("error resize pty: %s",err)
			}
		}
	}()
	ch <- syscall.SIGWINCH
	
	oldState,err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Error("error make raw: %s",err)
		return log.NewError("error make raw: %s",err)
	}
	
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()),oldState) }()
	
	go func() { ReadStdout(ptmx) }()
//	err = ReadInput()
	
//	go func() { _,_ = io.Copy(ptmx, os.Stdin) }()
//	_,_ = io.Copy(os.Stdout,ptmx)
	 _,_ = io.Copy(ptmx, os.Stdin)

	
	return err
}

func ReadStdout(r io.Reader) error {
	reader := bufio.NewReader(r)
	for {
		output, err := reader.ReadString('\n')
		if err == io.EOF { break }
		if err != nil {
			log.Debug("read stdout error: %s",err)
			break
		}
		if false { fmt.Fprintf(os.Stdout,output) }
		log.Debug(output)
	}
	return nil
}
//
//func ReadInput() error {
//	reader := bufio.NewReader(os.Stdin)
//	for {
//		input, err := reader.ReadString('\n')
//		if err != nil {
//			log.Debug("read error: %s",err)
//		}	
//		log.Debug(input)
//	}	
//	return nil
//}

	
//	log.Debug("exec cmd %s %s",cmd.Path,strings.Join(cmd.Args[1:]," "))
//
//
//    signals := make(chan os.Signal, 1)
//    signal.Notify(signals, os.Interrupt)
//    go func() {
//        sig := <-signals
//        log.Notice("%s",sig)
//    }()
//
//
//	var err error
//    var stdin, stdout, stderr bytes.Buffer
//
//	cmd.Stdout = &stdout
//	cmd.Stderr = &stderr
//	cmd.Stdin = &stdin
//	
//	stdout,err := cmd.StdoutPipe()
//	if err != nil { log.Error("error stdout pipe: %s",err) }
//	stderr,err := cmd.StderrPipe()
//	if err != nil { log.Error("error stderr pipe: %s",err) }
//	
//	go ReadStdout( stdout )
//	go ReadStderr( stderr )
//	
//	go ReadInput()
//	
//	err = cmd.Start()
//	
//	
//	
//	if err != nil { log.Error("error start cmd %s: %s",cmd.Path,err) }
//	
//	
//	
//	
//	err = cmd.Wait()
//	if err != nil { log.Error("error wait cmd %s: %s",cmd.Path,err) }
//	
//	
//	
//	
//	outStr, errStr := string(stdout.Bytes()),string(stderr.Bytes())
//	log.Debug(stdout)
//	log.Debug("")
//	log.Debug(stderr)
//	return err
//}
//
//
//
//func ReadStderr(r io.Reader) error {
//	reader  := bufio.NewReader(r)
//	for {
//		output, err := reader.ReadString('\n')
//		if err == io.EOF { break }
//		if err != nil {
//			log.Debug("read stderr error: %s",err)
//		}
//		fmt.Fprintf(os.Stderr,output)
//	}
//	return nil
//}
//
//func ReadInput() {
//	reader := bufio.NewReader(os.Stdin)
//	for {
//		input, err := reader.ReadString('\n')
//		if err != nil {
//			log.Debug("read error: %s",err)
//		}	
//		log.Debug(input)
//	}	
//}
