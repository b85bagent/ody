package exporter

import (
	model "Agent/model"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

var (
	cfg = model.Config{}
)

type MyExporter struct {
	mutex           sync.Mutex
	upMetric        prometheus.Gauge
	cpuUsageMetrics map[string]prometheus.Gauge
}

func (e *MyExporter) Describe(ch chan<- *prometheus.Desc) {
	e.upMetric.Describe(ch)
	for _, m := range e.cpuUsageMetrics {
		m.Describe(ch)
	}
}

func (e *MyExporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	log.Println("Collecting metrics...")

	for _, target := range cfg.Targets {
		data := getDataFromYAML(target)
		log.Println(data)
		e.upMetric.Set(data.Up)

		if data.CPUUsage.Enabled {
			e.cpuUsageMetrics[target].Set(data.CPUUsage.CPU1)
			e.cpuUsageMetrics[target+"-cpu2"].Set(data.CPUUsage.CPU2)
		}
	}

	e.upMetric.Collect(ch)
	for _, m := range e.cpuUsageMetrics {
		m.Collect(ch)
	}
}

func RunExporter() {
	yamlData, err := ioutil.ReadFile("snmp.yaml")
	if err != nil {
		log.Fatal("Failed to read YAML configuration:", err)
	}

	err = yaml.Unmarshal(yamlData, &cfg)
	if err != nil {
		log.Fatal("Failed to parse YAML configuration:", err)
	}

	exporter := &MyExporter{
		upMetric:        prometheus.NewGauge(prometheus.GaugeOpts{Name: "my_up_metric", Help: "Example exporter up metric"}),
		cpuUsageMetrics: make(map[string]prometheus.Gauge),
	}

	for _, target := range cfg.Targets {
		exporter.cpuUsageMetrics[target] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: fmt.Sprintf("snmp_%s", target),
			Help: fmt.Sprintf("snmp for %s", target),
		})
		exporter.cpuUsageMetrics[target+"-cpu2"] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: fmt.Sprintf("snmp_%s_cpu2", target),
			Help: fmt.Sprintf("cpu2 usage for %s", target),
		})
	}
	fmt.Println("exporter: ", exporter.cpuUsageMetrics)

	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Exporter listening on :8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getDataFromYAML(target string) model.Data {
	for _, metric := range cfg.Metrics {
		if metric.Name == target {
			return model.Data{
				Up: metric.CPUUsage.CPU1,
				CPUUsage: model.CPUUsageData{
					Enabled: true,
					CPU1:    metric.CPUUsage.CPU1,
					CPU2:    metric.CPUUsage.CPU2,
				},
			}
		}
	}

	return model.Data{}
}
