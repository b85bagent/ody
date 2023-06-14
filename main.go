package main

import (
	"Agent/pkg"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	bec "Agent/blackbox_exporter/config"
	bep "Agent/blackbox_exporter/prober"

	logger "github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 載入Server
	pkg.AutoLoader()

	target := "https://tw.yahoo.com"
	registry := prometheus.NewPedanticRegistry()
	t := bec.HTTPProbe{

		Method:             "GET",
		IPProtocolFallback: true,
		// IPProtocolFallback: true,
		// Compression:        "gzip",
		// Headers: map[string]string{
		// 	"Accept-Encoding": "*",
		// },
	}

	result := bep.ProbeHTTP(context.Background(), target,
		bec.Module{Timeout: time.Second, HTTP: t}, registry, logger.NewNopLogger())

	registry.MustRegister(NewHostMonitor())

	metrics, err := registry.Gather()
	if err != nil {
		log.Printf("Could not gather metrics: %v", err)
	}

	// 打印收集到的指标数据
	for _, mf := range metrics {
		// fmt.Printf("Metric: %s, %.1f\n", *mf.Name, *mf.Metric[0].Gauge.Value)
		for _, m := range mf.Metric {
			fmt.Printf("Metric: %s, Value: %v\n", *mf.Name, *m.Gauge.Value)
		}

	}

	if result {
		fmt.Println("Probe succeeded")
	} else {
		fmt.Println("Probe failed")
	}

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	http.ListenAndServe(":9090", nil)

}

type HostMonitor struct {
	cpuDesc    *prometheus.Desc
	memDesc    *prometheus.Desc
	ioDesc     *prometheus.Desc
	labelVaues []string
}

//创建结构体及对应的指标信息
func NewHostMonitor() *HostMonitor {
	return &HostMonitor{
		cpuDesc: prometheus.NewDesc(
			"host_cpu",
			"get host cpu",
			//动态标签key列表
			[]string{"instance_id", "instance_name"},
			//静态标签
			prometheus.Labels{"module": "cpu"},
		),
		memDesc: prometheus.NewDesc(
			"host_mem",
			"get host mem",
			//动态标签key列表
			[]string{"instance_id", "instance_name"},
			//静态标签
			prometheus.Labels{"module": "mem"},
		),
		ioDesc: prometheus.NewDesc(
			"host_io",
			"get host io",
			//动态标签key列表
			[]string{"instance_id", "instance_name"},
			//静态标签
			prometheus.Labels{"module": "io"},
		),
		labelVaues: []string{"myhost", "Lex"},
	}
}

//实现Describe接口，传递指标描述符到channel
func (h *HostMonitor) Describe(ch chan<- *prometheus.Desc) {
	ch <- h.cpuDesc
	ch <- h.memDesc
	ch <- h.ioDesc
}

//实现collect接口，将执行抓取函数并返回数据
func (h *HostMonitor) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(h.cpuDesc, prometheus.GaugeValue, 70, h.labelVaues...)
	ch <- prometheus.MustNewConstMetric(h.memDesc, prometheus.GaugeValue, 30, h.labelVaues...)
	ch <- prometheus.MustNewConstMetric(h.ioDesc, prometheus.GaugeValue, 90, h.labelVaues...)
}
