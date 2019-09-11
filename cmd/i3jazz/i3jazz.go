package main

import (
	"fmt"
	"os"
	"time"

	"github.com/erikh/i3bar"
)

func main() {
	ch := make(chan i3bar.StatusLine)
	go func() {
		for {
			ch <- i3bar.StatusLine{&i3bar.Block{FullText: "text"}, &i3bar.Block{FullText: fmt.Sprintf("%v", time.Now())}}
			time.Sleep(time.Second)
		}
	}()
	if err := i3bar.Encode(os.Stdout, &i3bar.Header{Version: 1}, ch); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
