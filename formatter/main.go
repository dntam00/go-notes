package main

import "fmt"
import "github.com/marmotedu/errors"

type dummyForFormatter struct {
}

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

func (d *dummyForFormatter) Format(s fmt.State, verb rune) {
	b := []byte{'k', 'a', 'i', 'x', 'i', 'n'}
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			_, _ = s.Write(b)
		case s.Flag('#'):
			// do something
		default:
			// do something
		}
	case 's':
		// do something
	}
}

func main() {
	//d := &dummyForFormatter{}
	//d.Format(newCustomState(), 'v')

	err()
}

func err() {
	err2 := chain1()
	if err2 != nil {
		fmt.Println(fmt.Sprintf("%+v", err2))
	}
}

func chain1() error {
	return chain2()
}

func chain2() error {
	return errors.WithCode(100, "dummy error")
}
