package ariago

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func Aria(url, outFile string) {
	aria := exec.Command("aria2c", "--stderr=true", "-c", "-s16", "-j16", "-x16", "-k1M", "-o", outFile, url)

	stderr, err := aria.StderrPipe()
	if err != nil {
		log.Println(err)
	}

	ch := make(chan string)
	done := make(chan interface{}, 1)

	go PipeRead(ch, done, stderr)

	aria.Start()

Loop:
	for {
		select {
		case s := <-ch:
			fmt.Print(s)
		case <-done:
			break Loop
		}
	}

	aria.Wait()
}

func PipeRead(ch chan string, done chan interface{}, r io.Reader) {
	rdr := bufio.NewReader(r)
	for {
		s, err := rdr.ReadString('\n')
		if err != nil {
			break
		}
		ch <- s
	}
	done <- nil
}
