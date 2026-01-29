package simpbar

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type SimpbarOption func(s *Simpbar)

func (s *Simpbar) SetOptions(options ...SimpbarOption) {
	for _, op := range options {
		op(s)
	}
}

func SimpbarSetTail(tail string) SimpbarOption {
	return func(s *Simpbar) {
		s.Tail = tail
	}
}

func SimpbarSetHead(head string) SimpbarOption {
	return func(s *Simpbar) {
		s.Head = head
	}
}
func SimpbarSetPref(pref string) SimpbarOption {
	return func(s *Simpbar) {
		s.Pref = pref
	}
}
func SimpbarSetSuff(suff string) SimpbarOption {
	return func(s *Simpbar) {
		s.Suff = suff
	}
}

type Simpbar struct {
	Width int

	Head string
	Tail string

	Pref string
	Suff string

	top  int
	base int

	since time.Time
	mu    sync.Mutex
}

func SimpbarInt(top int) Simpbar {
	return Simpbar{
		Width: 25,

		Head: "›",
		Tail: "-",

		Pref: "[",
		Suff: "]",

		top:  top,
		base: 0,
	}
}

func (s *Simpbar) Draw() {
	perc := float64(s.base) / float64(s.top)
	tail := int(perc * float64(s.Width))
	diff := s.Width - tail
	head := ""

	if diff == 0 {
		head = s.Tail
	} else {
		head = s.Head
	}

	bar := s.Pref + strings.Repeat(s.Tail, tail) + head + strings.Repeat(" ", diff) + s.Suff
	since := time.Since(s.since)
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("\r%s %s %.02f%% in %s.", timestamp, bar, perc*100, since.Round(time.Second))
}

func (s *Simpbar) Add(value int) {
	s.mu.Lock()
	s.base = min(s.base+value, s.top)

	if s.since.IsZero() {
		s.since = time.Now()
	}

	s.Draw()
	if s.base == s.top {
		fmt.Println()
	}
	s.mu.Unlock()
}
