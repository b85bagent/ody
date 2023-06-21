package exporter

import (
	"context"
	"log"
	"time"

	bec "Agent/blackbox_exporter/config"
	bep "Agent/blackbox_exporter/prober"

	logger "github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
)

func ProbeHttp(data map[string]interface{}, target string) (resultData map[string]interface{}) {

	log.Println("Probing http Start: ", target)

	registry := prometheus.NewPedanticRegistry()

	t := bec.HTTPProbe{
		IPProtocolFallback: true,
		// IPProtocolFallback: true,
		// Compression:        "gzip",
		// Headers: map[string]string{
		// 	"Accept-Encoding": "*",
		// },
	}

	timeout := timeOutSetting()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := bep.ProbeHTTP(ctx, target,
		bec.Module{Timeout: time.Microsecond, HTTP: t}, registry, logger.NewNopLogger())

	registry.MustRegister(NewHostMonitor())

	metrics, err := registry.Gather()
	if err != nil {
		log.Printf("Could not gather metrics: %v", err)
	}

	// 處理收集到的數據資料，寫入map內
	for _, mf := range metrics {
		r := make(map[string]interface{})
		nested := make(map[string]interface{})
		for i, m := range mf.Metric {
			if len(mf.Metric[i].Label) != 0 {
				name := *mf.Name
				if name == "probe_ssl_last_chain_info" {
					data[*mf.Name] = m.Gauge.Value
					continue
				}

				for _, v := range mf.Metric[i].Label {
					// labelName := *v.Name
					labelValue := *v.Value

					nested[labelValue] = m.Gauge.Value
					// r[labelName+"_"+fmt.Sprint(i)] = nested

				}

				r[*mf.Metric[i].Label[0].Name] = nested

				data[name] = r
			} else {
				data[*mf.Name] = m.Gauge.Value
			}
		}
	}

	if result {
		data["result"] = "Success"
		return data
	}

	data["result"] = "Failed"

	return data
}

func ProbeHttpPOST(data map[string]interface{}, target string) (resultData map[string]interface{}) {

	log.Println("Probing http Start: ", target)

	registry := prometheus.NewPedanticRegistry()

	t := bec.HTTPProbe{
		IPProtocolFallback: true,
		Method:             "POST",
		// IPProtocolFallback: true,
		// Compression:        "gzip",
		// Headers: map[string]string{
		// 	"Accept-Encoding": "*",
		// },
	}

	timeout := timeOutSetting()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := bep.ProbeHTTP(ctx, target,
		bec.Module{Timeout: time.Microsecond, HTTP: t}, registry, logger.NewNopLogger())

	registry.MustRegister(NewHostMonitor())

	metrics, err := registry.Gather()
	if err != nil {
		log.Printf("Could not gather metrics: %v", err)
	}

	// log.Println("metrics: ", metrics)

	// 處理收集到的數據資料，寫入map內
	for _, mf := range metrics {
		r := make(map[string]interface{})
		nested := make(map[string]interface{})
		for i, m := range mf.Metric {
			if len(mf.Metric[i].Label) != 0 {
				name := *mf.Name
				if name == "probe_ssl_last_chain_info" {
					data[*mf.Name] = m.Gauge.Value
					continue
				}

				for _, v := range mf.Metric[i].Label {
					// labelName := *v.Name
					labelValue := *v.Value

					nested[labelValue] = m.Gauge.Value
					// r[labelName+"_"+fmt.Sprint(i)] = nested

				}

				r[*mf.Metric[i].Label[0].Name] = nested

				data[name] = r
			} else {
				data[*mf.Name] = m.Gauge.Value
			}
		}
	}

	if result {
		data["result"] = "Success"
		return data
	}

	data["result"] = "Failed"

	return data
}
