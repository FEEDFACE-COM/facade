

package main

import (
    "os"
    "io"
    "bufio"
    log "./log"
    facade "./facade"
)


const DEBUG_READ = false

type Scanner struct {
    reader *bufio.Reader   
}

func NewScanner() *Scanner {
    ret := &Scanner{}
    ret.reader = bufio.NewReader(os.Stdin)
    return ret
}


func (scanner *Scanner) ScanText(bufChan chan facade.BufferItem) {
	var rem []byte = []byte{}
	var tmp []byte

    const BUFFER_SIZE = 1024
	var buf []byte = make([]byte, BUFFER_SIZE)
	
	for {
        n,err := scanner.reader.Read(buf)
		if err == io.EOF { 
    		if DEBUG_READ { log.Debug("read end of file") }
    		break 
        }
		if err != nil {
			log.Error("read stdin error: %s",err)
			break
		}
        if DEBUG_READ { log.Debug("read %d byte:\n%s",n,log.Dump(buf,n,0)) }

		tmp = append(rem, buf[:n] ... )
		rem, err = facade.ProcessRaw(tmp, bufChan)
		if err != nil {
            log.Error("process error: %s",err)    		
        }
	
    }    
}

