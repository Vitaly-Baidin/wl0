package zapLog

import "go.uber.org/zap"

type Logger struct {
	Logger *zap.SugaredLogger
}

func (l *Logger) Debug(args ...any) {
	l.Logger.Debug(args)
}

func (l *Logger) Debugf(template string, args ...any) {
	l.Logger.Debugf(template, args)
}

func (l *Logger) Debugln(args ...any) {
	l.Logger.Debugln(args)
}

func (l *Logger) Info(args ...any) {
	l.Logger.Info(args)
}

func (l *Logger) Infof(template string, args ...any) {
	l.Logger.Infof(template, args)
}

func (l *Logger) Infoln(args ...any) {
	l.Logger.Infoln(args)
}

func (l *Logger) Warn(args ...any) {
	l.Logger.Warn(args)
}

func (l *Logger) Warnf(template string, args ...any) {
	l.Logger.Warnf(template, args)
}

func (l *Logger) Warnln(args ...any) {
	l.Logger.Warnln(args)
}

func (l *Logger) Error(args ...any) {
	l.Logger.Error(args)
}

func (l *Logger) Errorf(template string, args ...any) {
	l.Logger.Errorf(template, args)
}

func (l *Logger) Errorln(args ...any) {
	l.Logger.Errorln(args)
}

func (l *Logger) Panic(args ...any) {
	l.Logger.Panic(args)
}

func (l *Logger) Panicf(template string, args ...any) {
	l.Logger.Panicf(template, args)
}

func (l *Logger) Panicln(args ...any) {
	l.Logger.Panicln(args)
}

func (l *Logger) Fatal(args ...any) {
	l.Logger.Fatal(args)
}

func (l *Logger) Fatalf(template string, args ...any) {
	l.Logger.Fatalf(template, args)
}

func (l *Logger) Fatalln(args ...any) {
	l.Logger.Fatalln(args)
}
