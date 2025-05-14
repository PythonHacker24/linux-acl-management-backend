package utils

import (
	"log"
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/models"
)

var (
	Log *zap.Logger
)

/* initializes the zap logger and provides global logging */
func InitLogger(isProduction bool) {
	var encoder zapcore.Encoder
	var writeSyncer zapcore.WriteSyncer
	var logLevel zapcore.Level

	/* check if the logging level is production */
	if isProduction {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		logLevel = zapcore.InfoLevel
		writeSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    100, // MB
			MaxBackups: 5,
			MaxAge:     30, // days
			Compress:   true,
		})
	} else {

		/* development level logging - configured for debug */

		cfg := zap.NewDevelopmentEncoderConfig()
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(cfg)
		logLevel = zapcore.DebugLevel
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	/* allow global logging with zap.L() */
	zap.ReplaceGlobals(Log)

	log.Println("Initialized Zap Logger")
}

/* generate a new uuid */
func GenerateTxnID() string {
	return uuid.New().String()
}
