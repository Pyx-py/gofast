package utils

import (
	"os"
	"path"
	"time"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
)

func GetWriteSyncer(linkName, logPath string, day int, logInConsole bool) (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(logPath, "%Y-%m-%d.log"),
		zaprotatelogs.WithLinkName(linkName),
		zaprotatelogs.WithMaxAge(time.Duration(day)*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if logInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
