package ui

type UI interface {
	ReadLine() (string, error)
	Printf(format string, a ...any)
	Println(a ...any)
}
