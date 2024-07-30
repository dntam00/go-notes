package main

import (
	"errors"
	"fmt"
	"runtime"
)

func main() {
	_ = test()
}

func test() error {
	st := callers()
	caller := fmt.Sprintf("#%d", 0)
	f := Frame((*st)[0])
	caller = fmt.Sprintf("%s %s:%d (%s)",
		caller,
		f.file(),
		f.line(),
		f.name(),
	)
	fmt.Print(caller)
	return errors.New("error")
}

type stack []uintptr

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := Frame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

type Frame uintptr

type customState struct{}

func newCustomState() *customState {
	return &customState{}
}

func (c *customState) Write(b []byte) (n int, err error) {
	return fmt.Print(string(b))
}

func (c *customState) Width() (wid int, ok bool) {
	return 0, false
}

func (c *customState) Precision() (prec int, ok bool) {
	return 0, false
}

func (c *customState) Flag(ch int) bool {
	return ch == '+'
}

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name returns the name of this function, if known.
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}
