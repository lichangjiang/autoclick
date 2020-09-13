package main

import (
	"autoclick/cmd"
	//"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:    "autoclick",
	Hidden: true,
}

func init() {
	// Log as JSON instead of the default ASCII formatter.

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	/*file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.SetOutput(os.Stdout)
		log.Info("Failed to log to file")
	}*/
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

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
