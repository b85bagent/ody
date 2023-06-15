package exporter

import (
	"context"
	"fmt"
	"log"
	"time"

	"Agent/server"
	bec "Agent/blackbox_exporter/config"
	bep "Agent/blackbox_exporter/prober"

	logger "github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/config"
)

func ProbeHttp(data map[string]interface{}, target string) (resultData map[string]interface{}) {

	fmt.Println("Probing http Start: ", target)
	httpClient := config.HTTPClientConfig{
		// 設置重試次數為 2 次
	}

	registry := prometheus.NewPedanticRegistry()
	t := bec.HTTPProbe{
		IPProtocolFallback: true,
		HTTPClientConfig:   httpClient,
		// IPProtocolFallback: true,
		// Compression:        "gzip",
		// Headers: map[string]string{
		// 	"Accept-Encoding": "*",
		// },
	}

	os, ok := server.GetServerInstance().GetConst()["httpRetrySecond"]
	if !ok {
		fmt.Println("No OK")
		return
	}

	timeout := time.Duration(os.(int)) * time.Second // 設定超時時間 5秒 修改請見config
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := bep.ProbeHTTP(ctx, target,
		bec.Module{Timeout: time.Microsecond, HTTP: t}, registry, logger.NewNopLogger())

	registry.MustRegister(NewHostMonitor())

	metrics, err := registry.Gather()
	if err != nil {
		log.Printf("Could not gather metrics: %v", err)
	}

	// fmt.Println("metrics: ", metrics)

	// 打印收集到的指标数据
	for _, mf := range metrics {
		// fmt.Printf("Metric: %s, %.1f\n", *mf.Name, *mf.Metric[0].Gauge.Value)
		for i, m := range mf.Metric {
			if len(mf.Metric[i].Label) != 0 {

				name := *mf.Name

				for _, v := range mf.Metric[i].Label {
					name = name + "{" + *v.Name + ":" + *v.Value + "}"
				}

				data[name] = m.Gauge.Value
				// fmt.Printf("Metric: %s, Metric: %v, Value: %v\n", *mf.Name, mf.Metric[i].Label, m.Gauge)
				continue
			}

			data[*mf.Name] = m.Gauge.Value
			// fmt.Printf("Metric: %s,Value: %v\n", *mf.Name, m.Gauge)

		}

	}

	if result {
		fmt.Println("Probe succeeded")
	} else {
		fmt.Println("Probe failed")
	}

	return data
}

type HostMonitor struct{}

//创建结构体及对应的指标信息
func NewHostMonitor() *HostMonitor {
	return &HostMonitor{}
}

//实现Describe接口，传递指标描述符到channel
func (h *HostMonitor) Describe(ch chan<- *prometheus.Desc) {}

//实现collect接口，将执行抓取函数并返回数据
func (h *HostMonitor) Collect(ch chan<- prometheus.Metric) {}
