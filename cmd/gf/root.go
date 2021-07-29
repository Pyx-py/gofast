package gf

import (
	"os"

	"github.com/gookit/color"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gf",
	Short: "这是一款用来生初始化代码文件的终端工具",
	Long: `欢迎使用gf终端工具
  ————————   ——————————
 /  —————/   |  ——————/
/   \  ———   |  |————
\    \——\ \  |   ————|
 \——————  /  |  |
        \/   |——|
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gva.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.AddCommand(fileCmd)
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		if err != nil {
// 			color.Warn.Println(err)
// 			os.Exit(1)
// 		}

// 		// Search config in home directory with name ".gva" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigName(".gva")
// 	}

// 	viper.AutomaticEnv() // read in environment variables that match

// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		color.Warn.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }
