package main

import (
	"autoclick/cmd"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:    "autoclick",
	Hidden: true,
}

func init() {
	logName := "autoclick_server.log"

	f, _ := os.OpenFile(logName, os.O_WRONLY|os.O_CREATE, 0755)
	mw := io.MultiWriter(f, os.Stdout)
	log.SetOutput(mw)
	log.SetFormatter(&log.JSONFormatter{})
	rootCmd.AddCommand(cmd.ServerCmd)
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{
			"launch": "error",
			"error":  err,
		}).Fatal("launch fail")
	}
}
