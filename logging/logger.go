package logging

type Logger interface {
	DebugLogger
	InfoLogger
	WarnLogger
	ErrorLogger
	PanicLogger
	FatalLogger
}

type DebugLogger interface {
	Debug(args ...any)
	Debugf(template string, args ...any)
	Debugln(args ...any)
}

type InfoLogger interface {
	Info(args ...any)
	Infof(template string, args ...any)
	Infoln(args ...any)
}

type WarnLogger interface {
	Warn(args ...any)
	Warnf(template string, args ...any)
	Warnln(args ...any)
}

type ErrorLogger interface {
	Error(args ...any)
	Errorf(template string, args ...any)
	Errorln(args ...any)
}

type PanicLogger interface {
	Panic(args ...any)
	Panicf(template string, args ...any)
	Panicln(args ...any)
}

type FatalLogger interface {
	Fatal(args ...any)
	Fatalf(template string, args ...any)
	Fatalln(args ...any)
}
