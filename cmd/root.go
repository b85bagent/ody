package cmd

import (
	"Agent/pkg"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func Run() {
	var snmpFile string
	var configFile string

	rootCmd := &cobra.Command{
		Use:   "Agent",
		Short: "Agent Application",
		Long:  "Agent using Cli for setting yaml configuration",
		Run: func(cmd *cobra.Command, args []string) {
			// 在這裡添加你的應用程式邏輯
			pkg.AutoLoader(configFile, snmpFile)
		},
	}

	// 添加 YAML 環境檔案選項
	rootCmd.PersistentFlags().StringVarP(&snmpFile, "snmp", "s", "snmp.yaml", "指定要載入的 SNMP YAML 檔案")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "指定要載入的設定 YAML 檔案")

	// 在執行根命令之前的預處理邏輯
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		log.Printf("使用的 config 檔案: %s / 使用的 SNMP 檔案: %s\n", configFile, snmpFile)
		// 檢查是否指定了自訂的 SNMP YAML 檔案
	}

	// 執行根命令
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
