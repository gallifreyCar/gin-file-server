package m_logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
)

func InitZapLogger(fileName, prefix string) (logger *zap.Logger, err error, close func()) {
	filePath := path.Join("..", "target", "log", fileName)
	logFile, closeLogFile, err := zap.Open(filePath)
	if err != nil {
		panic(err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		logFile,
		zap.NewAtomicLevel(),
	)
	prefixOption := zap.Fields(zap.String("func", prefix))
	logger = zap.New(core, prefixOption, zap.AddCaller())
	return logger, err, closeLogFile
}
