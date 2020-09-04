package main

import (
	facade "./facade"
	log "./log"
	"bufio"
	"fmt"
	"io"
	"os"
)

const DEBUG_SCAN = false
const DEBUG_SCAN_DUMP = false

type Scanner struct {
	reader     *bufio.Reader
	bufferSize uint
}

func NewScanner() *Scanner {
	ret := &Scanner{
		reader:     bufio.NewReader(os.Stdin),
		bufferSize: TEXT_BUFFER_SIZE,
	}
	return ret
}

func (scanner *Scanner) ScanText(bufChan chan facade.TextSeq) {
	var rem []byte = []byte{}
	var tmp []byte

	var buf []byte = make([]byte, scanner.bufferSize)

	log.Info("%s read text from stdin", scanner.Desc())

	for {
		n, err := scanner.reader.Read(buf)
		if err == io.EOF {
			if DEBUG_SCAN {
				log.Debug("%s read end of file", scanner.Desc())
			}
			break
		}
		if err != nil {
			log.Error("%s read stdin error: %s", scanner.Desc(), err)
			break
		}
		if DEBUG_SCAN_DUMP {
			log.Debug("%s read %d byte:\n%s", scanner.Desc(), n, log.Dump(buf, n, 0))
		} else if DEBUG_SCAN {
			log.Debug("%s read %d byte", scanner.Desc(), n)
		}

		tmp = append(rem, buf[:n]...)
		rem, err = facade.ProcessRaw(tmp, bufChan)
		if err != nil {
			log.Error("%s process error: %s", scanner.Desc(), err)
		}

	}
}

func (scanner *Scanner) Desc() string {
	return fmt.Sprintf("scanner[%d]", scanner.bufferSize)
}
