package logs

import (
	//Inbuild packages
	"fmt"
	"io"
	"os"

	//Third-party packages
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func Log() *logrus.Logger {
	f, err := os.OpenFile("log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
		panic(err)
	}

	log := &logrus.Logger{
		// Log into f file handler and on os.Stdout
		Out:   io.MultiWriter(f, os.Stdout),
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "02-01-2006 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}
	return log
}
