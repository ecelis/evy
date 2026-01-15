package main

import (
	"github.com/nsf/termbox-go"
	"log"
)

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		termbox.Flush()

		switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyEsc {
					return
				}
			case termbox.EventError:
				log.Fatal(ev.Err) 
		}  
	} 
}
