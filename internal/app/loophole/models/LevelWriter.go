package models

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//Implementation of io.Writer which writes to zerolog with a configurable level
//and identifying prefix to prevent "???" level messages
type LevelWriter struct {
	Level         zerolog.Level
	MessagePrefix string
}

func (l LevelWriter) Write(p []byte) (n int, err error) {
	log.WithLevel(l.Level).Msg(fmt.Sprintf("%s %s", l.MessagePrefix, string(p)))
	return len(p), nil
}
