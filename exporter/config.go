package exporter

import (
	"Agent/server"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type HostMonitor struct{}

//創建結構體對應指標
func NewHostMonitor() *HostMonitor {
	return &HostMonitor{}
}

//實現Describe接口，傳遞指標到channel
func (h *HostMonitor) Describe(ch chan<- *prometheus.Desc) {}

//實現collect接口，執行抓取函數返回數據
func (h *HostMonitor) Collect(ch chan<- prometheus.Metric) {}

// 確認module類型，給予不同的Probe
func CheckModule(module string, data map[string]interface{}, target string) (resultData map[string]interface{}, err error) {
	switch module {
	case "http_2xx":
		resultData = ProbeHttp(data, target)
	case "http_post_2xx":
		resultData = ProbeHttpPOST(data, target)
	case "dns":
		resultData = ProbeDns(data, target)
	case "tcp_connect":
		resultData = ProbeTcp(data, target)
	case "icmp":
		resultData = ProbeIcmp(data, target)
	case "grpc":
		resultData = ProbeGrpc(data, target)
	}

	return resultData, nil
}

//設定監控每次probe的timeout時間
func timeOutSetting() time.Duration {
	timeoutSetting, ok := server.GetServerInstance().GetConst()["httpRetrySecond"]
	if !ok {
		fmt.Println("timeoutSetting Get failed")
		return 0
	}

	timeout := time.Duration(timeoutSetting.(int)) * time.Second // 設定超時時間 5秒 修改請見config.yaml

	return timeout
}
