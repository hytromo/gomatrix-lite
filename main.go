package main

import (
	"fmt"
	"log"

	"os"

	"github.com/gdamore/tcell/v2"
)

const VERSION = "0.0.1"

func eventLoop(xmax *int, ymax *int, waitTimeMs *int64, _s *tcell.Screen) {
	s := *_s

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			*xmax, *ymax = s.Size()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'Q' || ev.Rune() == 'q' {
				s.Fini()
				os.Exit(0)
			} else if ev.Rune() >= '0' && ev.Rune() <= '9' {
				(*waitTimeMs) = 95 - int64((ev.Rune()-'0')*10)
			}
		}
	}
}

func main() {
	config := ParseArgs()

	if config.showVersion {
		fmt.Println(VERSION)
		return
	}

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.EnablePaste()
	s.Clear()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	xmax, ymax := s.Size()
	var waitTime int64 = 35

	go Matrix(&xmax, &ymax, &waitTime, &config.colors, &s)
	eventLoop(&xmax, &ymax, &waitTime, &s)
}
