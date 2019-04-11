package cmd

import (
	"fmt"
	"github.com/insufficientchocolate/diplomat"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"path/filepath"
)

var (
	watch bool
	buildCmd = &cobra.Command{
		Use: "build",
		Short: "build",
		Run: func(cmd *cobra.Command, args []string) {
			var projectDir string
			if len(args) > 0 {
				projectDir = args[0]
			} else {
				var err error
				projectDir, err = os.Getwd()
				if err != nil {
					fmt.Println("cannot get current working directory", err)
					return
				}
			}
			outDir := filepath.Join(projectDir, "out")
			if watch {
				d, errChan := diplomat.NewDiplomatWatchDirectory(projectDir)
				go func() {
					for e := range errChan {
						log.Println("error:", e)
					}
				}()
				go d.Watch(outDir)
				quit := make(chan os.Signal, 1)
				signal.Notify(quit, os.Interrupt)
				<-quit
			} else {
				d, err := diplomat.NewDiplomatForDirectory(projectDir)
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
				if err := d.Output(outDir); err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
			}

		},
	}
)

func init() {
	buildCmd.Flags().BoolVar(&watch, "watch",false,"watch changes")
}
