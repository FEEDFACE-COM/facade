

package main

import (
    "os"
    "bufio"
    log "./log"
    render "./render"
)


type Scanner struct {
    scanner *bufio.Scanner   
}

func NewScanner() *Scanner {
    ret := &Scanner{}
    ret.scanner = bufio.NewScanner(os.Stdin)
    return ret
}


func (scanner *Scanner) ScanText(textChan chan render.Text) {
    for scanner.scanner.Scan() {
        text := scanner.scanner.Text()
//        log.Debug("scan  %s",text)
        textChan <- render.Text(text)
    }    
    err := scanner.scanner.Err()
    if err != nil {
        log.Error("error scanning: %s",err)
    }
}

