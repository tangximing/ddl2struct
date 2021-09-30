package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tangximing/ddl2struct/gen"
)

func Do() {
	rootCmd := getRootCmd()
	err := rootCmd.Execute()
	if err == nil {
		os.Exit(0)
	}

	os.Exit(1)
}

func getRootCmd() *cobra.Command {
	var sqlFile string
	var outDir string
	var packageName string

	rootCmd := &cobra.Command{
		Use:          "ddl2struct",
		Short:        "generate golang struct file from ddl",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return gen.Generate(sqlFile, outDir, packageName)
		},
	}
	rootCmd.Flags().StringVarP(&sqlFile, "sql", "s", "", "ddl sql file path")
	rootCmd.Flags().StringVarP(&outDir, "dir", "d", "", "golang dir to generate")
	rootCmd.Flags().StringVarP(&packageName, "package", "p", "", "golang package to generate")

	return rootCmd
}
