package debug

import (
	"os"

	"github.com/spf13/cobra"
)

var backtraceCmd = &cobra.Command{
	Use:     "bt",
	Short:   "打印调用栈信息",
	Aliases: []string{"backtrace"},
	Annotations: map[string]string{
		cmdGroupAnnotation: cmdGroupInfo,
	},
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}

func init() {
	debugRootCmd.AddCommand(backtraceCmd)
}
