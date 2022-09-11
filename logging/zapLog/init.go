package zapLog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

func InitializeLogger(path ...string) (*Logger, error) {
	fileEncoder := createFileEncode()
	consoleEncoder := createConsoleEncode()
	fp := filepath.Join(path...)

	allWriter, err := createWriter(fp, "/all.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create writer all logs: %w\n", err)
	}
	errorWriter, err := createWriter(fp, "/err.log")
	if err != nil {
		return nil, fmt.Errorf("failed to create writer err logs: %w\n", err)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, allWriter, zapcore.DebugLevel),
		zapcore.NewCore(fileEncoder, errorWriter, zapcore.ErrorLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)
	sugaredLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()

	return &Logger{
		Logger: sugaredLogger,
	}, nil
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

func createWriter(folder, file string) (zapcore.WriteSyncer, error) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w\n", err)
	}

	logFile, err := os.OpenFile(folder+file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open folder: %w\n", err)
	}

	return zapcore.AddSync(logFile), nil
}
