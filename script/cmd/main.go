package main

import (
	"deposit/cmd/randomtransfer"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	
	RootCmd = &cobra.Command{
		Use:   "example",
		Short: "An example cobra program",
		Long:  "This is a simple cobra program",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello from root cmd")
		},
	}
)

func init() {
	RootCmd.AddCommand(randomtransfer.Transfer)
	
}
func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
		PrettyPrint:      true,
	})
	logrus.SetReportCaller(true)
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
