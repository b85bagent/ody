package exporter

import (
	"context"
	"fmt"
	"log"
	"time"

	bec "Agent/blackbox_exporter/config"
	bep "Agent/blackbox_exporter/prober"

	logger "github.com/go-kit/log"

	"github.com/prometheus/client_golang/prometheus"
)

func ProbeGrpc(data map[string]interface{}, target string) (resultData map[string]interface{}) {

	fmt.Println("Probing Grpc Start: ", target)

	registry := prometheus.NewPedanticRegistry()

	t := bec.GRPCProbe{
		IPProtocolFallback: true,
	}

	timeout := timeOutSetting()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := bep.ProbeGRPC(ctx, target,
		bec.Module{Timeout: time.Microsecond, GRPC: t}, registry, logger.NewNopLogger())

	registry.MustRegister(NewHostMonitor())

	metrics, err := registry.Gather()
	if err != nil {
		log.Printf("Could not gather metrics: %v", err)
	}

	fmt.Println("metrics: ", metrics)

	// 處理收集到的數據資料，寫入map內
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
		data["result"] = "GRPC_Probe success"
		return data
	}

	data["result"] = "GRPC_Probe failed"

	return data
}
