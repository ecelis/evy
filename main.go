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
	"log"

	"github.com/nsf/termbox-go"
)

type Editor struct {
	lines []string
	curL, curC int
	mode int
} 

const (
	ModeNormal = iota
	ModeInsert
) 

var edit = Editor{
	lines: []string{""},
	mode: ModeNormal,
}

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for l, line := range edit.lines {
		for c, char := range line {
			termbox.SetCell(c, l, char, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	termbox.SetCursor(edit.curC, edit.curL)
	termbox.Flush() 
}

func handleKey(ev termbox.Event) {
	if edit.mode == ModeInsert {
		handleInsertMode(ev) 
	} else {
		handleNormalMode(ev)
	} 
} 

func handleInsertMode(ev termbox.Event) {
	switch ev.Key {
		case termbox.KeyEsc:
			edit.mode = ModeNormal
		case termbox.KeyArrowLeft:
			if edit.curC > 0  {edit.curC--}
		case termbox.KeyArrowRight:
			if edit.curC < len(edit.lines[edit.curL] ) {edit.curC++}
		case termbox.KeyArrowUp:
			if edit.curL > 0 {edit.curL--}
		case termbox.KeyArrowDown:
			if edit.curL < len(edit.lines)-1 {edit.curL++}
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
		default:
			if ev.Ch != 0 {
				line := edit.lines[edit.curL]
				edit.lines[edit.curL] = line[:edit.curC] + string(ev.Ch) + line[edit.curC:]
				edit.curC++
			} 
	} 
}

func handleNormalMode(ev termbox.Event) {
	switch ev.Ch {
		case 'i':  
			edit.mode = ModeInsert
		case 'h':
			if edit.curC > 0 {edit.curC--} 
		case 'j':
			if edit.curL < len(edit.lines) - 1 {edit.curL++} 
		case 'k':
			if edit.curL > 0 {edit.curL--} 
		case 'l':
			if edit.curC < len(edit.lines[edit.curL]) {edit.curC++}  
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
	}
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
				handleKey(ev) 
			case termbox.EventError:
				log.Fatal(ev.Err) 
		}  
	} 
}
