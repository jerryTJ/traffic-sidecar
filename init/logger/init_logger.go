package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger_file string

func Init(file string, level zerolog.Level) {
	Logger_file = file
	zerolog.SetGlobalLevel(level)

}
func Info(content string) {
	fileHandler, err := os.OpenFile(Logger_file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fileHandler.Close()
	logger := zerolog.New(fileHandler).With().Timestamp().Logger()
	logger.Info().Msg(content)
}
