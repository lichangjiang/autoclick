package main

import (
	"autoclick/cmd"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:    "autoclick",
	Hidden: true,
}

func init() {
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
