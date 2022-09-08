package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
)

var Logger *zap.SugaredLogger

func InitializeLogger(path ...string) {
	fileEncoder := createFileEncode()
	consoleEncoder := createConsoleEncode()
	fp := filepath.Join(path...)

	allWriter := createWriter(fp, "/all.log")
	errorWriter := createWriter(fp, "/error.log")

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, allWriter, zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, errorWriter, zapcore.ErrorLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
}

func createConsoleEncode() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(config)
}

func createFileEncode() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(config)
}

func createWriter(folder, file string) zapcore.WriteSyncer {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to create folder %s", err.Error())
	}

	logFile, err := os.OpenFile(folder+file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to create folder %s\n", err.Error())
	}

	return zapcore.AddSync(logFile)
}
