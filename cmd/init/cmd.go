package init

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "init",
	Short: "init Workorderdeployment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init ...... wait a moment")
	},
}

func init() {

}
