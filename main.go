package main

import (
	"log"

	"github.com/nsf/termbox-go"
)

type Editor struct {
	lines []string
	curX, curY int
} 

var edit = Editor{
	lines: []string{""},
} 

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for y, line := range edit.lines {
		for x, char := range line {
			termbox.SetCell(x, y, char, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	termbox.SetCursor(edit.curX, edit.curY)
	termbox.Flush() 
}

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	for {
		draw() 

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
