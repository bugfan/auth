package com

import (
	"os"

	"github.com/bugfan/logrus"
	"github.com/bugfan/to"
)

type EmptyWriter struct {
}

func (self *EmptyWriter) Write(p []byte) (n int, err error) {
	// nothing to do
	return len(p), nil
}

type LogrusWriter struct {
}

func (self *LogrusWriter) Write(p []byte) (n int, err error) {
	logrus.Infof(string(p))
	return len(p), nil
}
func init() {
	logrus.SetLevel(logrus.InfoLevel)
	if !to.Bool(os.Getenv("DEBUG")) {
		logrus.SetOutput(&EmptyWriter{})
	}
}
