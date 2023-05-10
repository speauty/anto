package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "anto",
	Short: "Anto is a subtitle translator, which based on cloud-service",
	Long: `
	Anto is a fast subtitle translator, you can use it to translate your srt files multiply.
website: github.com/speauty/anto
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("hello, my cobra")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
