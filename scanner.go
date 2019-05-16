

package main

import (
    "os"
    "bufio"
    log "./log"
    facade "./facade"
)


type Scanner struct {
    scanner *bufio.Scanner   
}

func NewScanner() *Scanner {
    ret := &Scanner{}
    ret.scanner = bufio.NewScanner(os.Stdin)
    return ret
}


func (scanner *Scanner) ScanText(bufChan chan facade.BufferItem) {
    var err error 
	var rem []byte = []byte{}
	var tmp []byte
    for scanner.scanner.Scan() {
        buf := scanner.scanner.Text()
		tmp = append(rem, buf ... )
		_, err = facade.ProcessRaw(tmp, bufChan)
		if err != nil {
            log.Error("process error: %s",err)    		
        }
    }    
    err = scanner.scanner.Err()
    if err != nil {
        log.Error("scanning error: %s",err)
    }
}

