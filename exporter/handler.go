package exporter

import (
	bec "Agent/blackbox_exporter/config"
	bep "Agent/blackbox_exporter/prober"
	"context"
	"errors"
	"log"

	logger "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// 確認module類型，給予不同的Probe
func CheckModuleAndDoProbe(module string, data map[string]interface{}, target string, sc *bec.SafeConfig) (resultData map[string]interface{}, err error) {

	result, err := comparisonConfigAndDoProbe(data, module, target, sc)
	if err != nil {
		log.Println("ReLoadConfig error: ", err)
		return nil, err
	}

	return result, nil
}

//比對yaml檔內容，並且Probe
func comparisonConfigAndDoProbe(data map[string]interface{}, m, target string, sc *bec.SafeConfig) (resultData map[string]interface{}, err error) {

	var e error
	logger := logger.NewNopLogger()

	level.Info(logger).Log("msg", "Reloaded config file")

	//comparisonConfig
	module, ok := sc.C.Modules[m]

	if !ok {
		level.Error(logger).Log("msg", "Module "+m+" not found")
		e = errors.New("Module " + m + " not found")

		return nil, e
	}

	prober, ok := Probers[module.Prober]

	if !ok {
		level.Error(logger).Log("msg", "Prober: "+module.Prober+"not found")
		e = errors.New("Prober: " + module.Prober + "not found")

		return nil, e
	}

	//doProbe
	result, errProbe := doProbe(data, module, prober, target)
	if errProbe != nil {
		level.Error(logger).Log("msg", "Probe failed", "err", errProbe)
		log.Println("Probe failed: ", errProbe)
		return nil, err
	}

	return result, nil
}

//Probe
func doProbe(data map[string]interface{}, module bec.Module, prober bep.ProbeFn, target string) (resultData map[string]interface{}, err error) {
	logger := logger.NewNopLogger()

	registry := prometheus.NewPedanticRegistry()

	timeout := timeOutSetting()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	success := prober(ctx, target, module, registry, logger)

	registry.MustRegister(NewHostMonitor())

	metrics, err := registry.Gather()
	if err != nil {
		log.Printf("Could not gather metrics: %v", err)
		return nil, err
	}

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

	if success {
		data["result"] = "Success"
		return data, nil
	}

	data["result"] = "Failed"

	return data, nil
}
