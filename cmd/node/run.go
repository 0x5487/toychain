package node

import (
	"fmt"
	stdlog "log"

	"github.com/nite-coder/blackbear/pkg/log"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "node",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		defer log.Flush()
		defer func() {
			if r := recover(); r != nil {
				// unknown error
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("unknown error: %v", r)
				}
				log.Err(err).Panicf("unknown error %v", err)
			}
		}()

		err := initialize()
		if err != nil {
			stdlog.Println(err.Error())
			log.Err(err).Panicf("main: toy chain initialize failed %v", err)
			return
		}
		defer db.Close()

		webServer.Run(":8080")
	},
}
