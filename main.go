package main

import (
	"clnotifications/clnotifications"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

func init() {
	logfilename := filepath.Join(os.Getenv("temp"), "clnotifications.log")
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   logfilename,
		MaxSize:    5, // megabytes
		MaxBackups: 5, // number of old log files to retain
		MaxAge:     1, // days
		Formatter: &log.TextFormatter{
			TimestampFormat: time.RFC822,
		},
	})

	if err != nil {
		log.Fatalf("Failed to initialize log file: %v", err)
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	log.AddHook(rotateFileHook)
}

func main() {
	log := log.WithFields(log.Fields{"event": "main"})
	fmt.Println("clnotifications v0.1.1")

	//err := clnotifications.GetKeys()
	err := clnotifications.ClearValues()
	if err != nil {
		fmt.Printf("Error during GetKeys: %v", err)
	}

}
