package cmd

import (
	"log"
	"os"
	"remote_write/pkg/autoload"

	"github.com/spf13/cobra"
)

func Run() {
	var configFile string

	rootCmd := &cobra.Command{
		Use:   "remote_write",
		Short: "remote_write Application",
		Long:  "remote_write using Cli for setting yaml configuration",
		Run: func(cmd *cobra.Command, args []string) {
			// 在這裡添加你的應用程式邏輯
			// autoload.AutoLoader(configFile, targetFile, blackboxFile)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "指定要載入的設定 Config YAML 檔案 (檔案請以config開頭)")

	// 在執行根命令之前的預處理邏輯
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		log.Printf("使用的 Config 檔案: %s\n", configFile)
		autoload.AutoLoader(configFile)
	}

	// 執行根命令
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
