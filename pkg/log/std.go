package log

import "fmt"

type StdLogger struct {
	l Logger
}

func (l Logger) Std() StdLogger {
	return StdLogger{l: l}
}

func (l StdLogger) Fatal(v ...any)                 { l.l.Fatal(fmt.Sprint(v...)) }
func (l StdLogger) Fatalf(format string, v ...any) { l.l.Fatal(fmt.Sprintf(format, v...)) }
func (l StdLogger) Fatalln(v ...any)               { l.l.Fatal(fmt.Sprint(v...)) }
func (l StdLogger) Panic(v ...any)                 { l.l.Panic(fmt.Sprint(v...)) }
func (l StdLogger) Panicf(format string, v ...any) { l.l.Panic(fmt.Sprintf(format, v...)) }
func (l StdLogger) Panicln(v ...any)               { l.l.Panic(fmt.Sprint(v...)) }
func (l StdLogger) Print(v ...any)                 { l.l.Info(fmt.Sprint(v...)) }
func (l StdLogger) Printf(format string, v ...any) { l.l.Info(fmt.Sprintf(format, v...)) }
func (l StdLogger) Println(v ...any)               { l.l.Info(fmt.Sprint(v...)) }
