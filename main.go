/**
* evy - Minimal command line text editor
*
* Copyright 2026 Ernesto Celis <ernesto@celisdelafuente.net>
*
* Permission is hereby granted, free of charge, to any person obtaining a copy of
* this software and associated documentation files (the “Software”), to deal in
* the Software without restriction, including without limitation the rights to
* use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
* of the Software, and to permit persons to whom the Software is furnished to do
* so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
* SOFTWARE.
**/

package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

type Editor struct {
	lines       []string
	curL, curC  int
	mode        int
	vOffset     int
	commandLine string
}

const (
	ModeNormal = iota
	ModeInsert
	ModeCommand
)

var edit = Editor{
	lines: []string{""},
	mode:  ModeNormal,
}

func draw() {
	edit.normalizeCursor()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()

	for l := 0; l < h-1; l++ {
		lineIdx := l + edit.vOffset
		if lineIdx < len(edit.lines) {
			for c, char := range edit.lines[lineIdx] {
				if c < w {
					termbox.SetCell(c, l, char, termbox.ColorWhite, termbox.ColorDefault)
				}
			}
		}
	}

	if edit.mode == ModeCommand {
		cmdStr := ":" + edit.commandLine
		for c, char := range cmdStr {
			termbox.SetCell(c, h-1, char, termbox.ColorDefault, termbox.ColorDefault)
		}
		termbox.SetCursor(len(cmdStr), h-1)
	} else {

		status := "-- NORMAL --"
		if edit.mode == ModeInsert {
			status = "-- INSERT --"
		}

		for c, char := range status {
			termbox.SetCell(c, h-1, char, termbox.ColorBlack, termbox.ColorWhite)
		}

		termbox.SetCursor(edit.curC, edit.curL)
	}
	termbox.Flush()
}

func executeCommand() {
	parts := strings.Split(edit.commandLine, " ")
	cmd := parts[0]

	switch cmd {
	case "q":
		termbox.Close()
		os.Exit(0)
	case "w":
		filename := ".swapfile"
		if len(parts) > 1 {
			filename = parts[1]
		}
		content := strings.Join(edit.lines, "\n")
		_ = os.WriteFile(filename, []byte(content), 0644)
	case "wq":
		executeCommand()
		termbox.Close()
		os.Exit(0)
	}
}

func handleKey(ev termbox.Event) {
	if edit.mode == ModeInsert {
		handleInsertMode(ev)
	} else if edit.mode == ModeCommand {
		handleCommandMode(ev)
	} else {
		handleNormalMode(ev)
	}
}

func handleInsertMode(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyEsc:
		edit.mode = ModeNormal
	case termbox.KeyArrowLeft:
		if edit.curC > 0 {
			edit.curC--
		}
	case termbox.KeyArrowRight:
		if edit.curC < len(edit.lines[edit.curL]) {
			edit.curC++
		}
	case termbox.KeyArrowUp:
		if edit.curL > 0 {
			edit.curL--
		}
	case termbox.KeyArrowDown:
		if edit.curL < len(edit.lines)-1 {
			edit.curL++
		}
	case termbox.KeyEnter:
		car := edit.lines[edit.curL][:edit.curC]
		cdr := edit.lines[edit.curL][edit.curC:]
		edit.lines[edit.curL] = car
		rest := append([]string{cdr}, edit.lines[edit.curL+1:]...)
		edit.lines = append(edit.lines[:edit.curL+1], rest...)
		edit.curL++
		edit.curC = 0
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if edit.curC > 0 {
			line := edit.lines[edit.curL]
			edit.lines[edit.curL] = line[:edit.curC-1] + line[edit.curC:]
			edit.curC--
		} else if edit.curL > 0 {
			prevLlen := len(edit.lines[edit.curL-1])
			edit.lines[edit.curL-1] += edit.lines[edit.curL]
			edit.lines = append(edit.lines[:edit.curL], edit.lines[edit.curL+1:]...)
			edit.curL--
			edit.curC = prevLlen
		}
	case termbox.KeySpace:
		line := edit.lines[edit.curL]
		edit.lines[edit.curL] = line[:edit.curC] + " " + line[edit.curC:]
		edit.curC++
	case termbox.KeyTab:
		//for i := 0; i < 4; i++ { 
			line := edit.lines[edit.curL]
			edit.lines[edit.curL] = line[:edit.curC] + "\t" + line[edit.curC:]
			edit.curC++
		//}
	default:
		if ev.Ch != 0 {
			line := edit.lines[edit.curL]
			edit.lines[edit.curL] = line[:edit.curC] + string(ev.Ch) + line[edit.curC:]
			edit.curC++
		}
	}
}

func handleNormalMode(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyArrowLeft:
		if edit.curC > 0 {
			edit.curC--
		}
	case termbox.KeyArrowDown:
		if edit.curL < len(edit.lines)-1 {
			edit.curL++
		}
	case termbox.KeyArrowUp:
		if edit.curL > 0 {
			edit.curL--
		}
	case termbox.KeyArrowRight:
		if edit.curC < len(edit.lines[edit.curL]) {
			edit.curC++
		}

	}
	switch ev.Ch {
	case 'i':
		edit.mode = ModeInsert
	case 'h':
		if edit.curC > 0 {
			edit.curC--
		}
	case 'j':
		if edit.curL < len(edit.lines)-1 {
			edit.curL++
		}
	case 'k':
		if edit.curL > 0 {
			edit.curL--
		}
	case 'l':
		if edit.curC < len(edit.lines[edit.curL]) {
			edit.curC++
		}
	case 'o':
		edit.lines = append(edit.lines[:edit.curL+1], append([]string{""}, edit.lines[edit.curL+1:]...)...)
		edit.curL++
		edit.curC = 0
		edit.mode = ModeInsert
	case 'x':
		line := edit.lines[edit.curL]
		if len(line) > 0 && edit.curC < len(line) {
			edit.lines[edit.curL] = line[:edit.curC] + line[edit.curC+1:]
		}
	case ':':
		edit.mode = ModeCommand
		edit.commandLine = ""
	}
}

func handleCommandMode(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyEsc:
		edit.mode = ModeNormal
	case termbox.KeyEnter:
		executeCommand()
		edit.mode = ModeNormal
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if len(edit.commandLine) > 0 {
			edit.commandLine = edit.commandLine[:len(edit.commandLine)-1]
		} else {
			edit.mode = ModeNormal
		}
	default:
		if ev.Ch != 0 {
			edit.commandLine += string(ev.Ch)
		}
	}
}

func (e *Editor) normalizeCursor() {
	if e.curL < 0 {
		e.curL = 0
	}
	if e.curL >= len(e.lines) {
		e.curL = len(e.lines) - 1
	}
	if e.curC < 0 {
		e.curC = 0
	}
	curLlen := len(e.lines[e.curL])
	if e.curC > curLlen {
		e.curC = curLlen
	}

	_, termHeight := termbox.Size()
	textHeight := termHeight - 1
	if e.curL >= e.vOffset+textHeight {
		e.vOffset = e.curL - textHeight + 1
	}
	if e.curL < e.vOffset {
		e.vOffset = e.curL
	}
}

func (e *Editor) readFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		e.lines = []string{""}
		return
	}
	defer file.Close() 

	var lines []string
	scanner := bufio.NewScanner(file) 
	for scanner.Scan() {
		lines = append(lines, scanner.Text()) 
	}
	if len(lines) == 0 {
		lines = append(lines, "") 
	}
	e.lines = lines
}   

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	if len(os.Args) > 1 {
		edit.readFile(os.Args[1]) 
	}

	for {
		draw()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			handleKey(ev)
		case termbox.EventError:
			log.Fatal(ev.Err)
		}
	}
}
