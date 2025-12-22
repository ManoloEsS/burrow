package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Console struct {
	in  *bufio.Reader
	out io.Writer
}

func NewConsole(r io.Reader, w io.Writer) *Console {
	return &Console{in: bufio.NewReader(r), out: w}
}

func (c *Console) ReadLine() (string, error) {
	s, err := c.in.ReadString('\n')
	return strings.TrimRight(s, "\r\n"), err
}

func (c *Console) Printf(format string, a ...any) {
	fmt.Fprintf(c.out, format, a...)
}

func (c *Console) Println(a ...any) {
	fmt.Fprintln(c.out, a...)
}
