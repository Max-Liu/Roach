package roach

type Logger interface {
	Warning(format string, v ...interface{})
	Informational(format string, v ...interface{})
	SetLogger(adaptername string, config string) error
}
