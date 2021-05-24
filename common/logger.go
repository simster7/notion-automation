package common

import log "github.com/sirupsen/logrus"

var logger *log.Entry

func InitLogger(job string) {
	if logger == nil {
		logger = log.WithFields(log.Fields{
			"job": job,
		})
	}
}

func GetLogger() *log.Entry {
	if logger == nil {
		log.Fatal("logger not initialized")
	}
	return logger
}
