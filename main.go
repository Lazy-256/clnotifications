package main

import (
	"clnotifications/clnotifications"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

func init() {
	exePath, err := os.Executable()
	if err != nil {
		exePath = os.Getenv("temp")
	}
	logfilename := filepath.Join(filepath.Dir(exePath), "clnotifications.log")
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   logfilename,
		MaxSize:    10, // megabytes
		MaxBackups: 1,  // number of old log files to retain
		MaxAge:     14, // days
		Level:      log.InfoLevel,
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

	flag_cleanup := flag.Bool("cleanup", true, "command to start cleaning up")
	flag.Parse()

	if !*flag_cleanup {
		err := clnotifications.GetKeys(log)
		if err != nil {
			fmt.Printf("Error during GetKeys: %v", err)
			log.Fatalf("Error during GetKeys: %v", err)
		}
	} else {
		log.Debugf("Startting cleanup process...")
		err := clnotifications.ClearValues(log)
		if err != nil {
			fmt.Printf("Error during ClearValues: %v", err)
			log.Fatalf("Error during ClearValues: %v", err)
		}
	}

}
