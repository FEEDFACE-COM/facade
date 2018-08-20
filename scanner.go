

package main

import (
    "os"
    "bufio"
    log "./log"
    conf "./conf"
)


type Scanner struct {
    scanner *bufio.Scanner   
}

func NewScanner() *Scanner {
    ret := &Scanner{}
    ret.scanner = bufio.NewScanner(os.Stdin)
    return ret
}


func (scanner *Scanner) ScanText(textChan chan conf.Text) {
    for scanner.scanner.Scan() {
        text := scanner.scanner.Text()
//        log.Debug("scan  %s",text)
        textChan <- conf.Text(text)
    }    
    err := scanner.scanner.Err()
    if err != nil {
        log.Error("error scanning: %s",err)
    }
}

