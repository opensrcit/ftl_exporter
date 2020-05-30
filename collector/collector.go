// Copyright 2020 Ivan Pushkin
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"flag"
	"fmt"
	"github.com/opensrcit/ftl_exporter/ftl_client"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"sync"
	"time"
)

var (
	factories      = make(map[string]func() (Collector, error))
	collectorState = make(map[string]*bool)
)

const (
	namespace       = "ftl"
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"ftl_exporter: Duration of a collector scrape.",
		[]string{"collector"}, nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"ftl_exporter: Whether a collector succeeded.",
		[]string{"collector"}, nil,
	)
)

func registerCollector(collector string, isDefaultEnabled bool, factory func() (Collector, error)) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector)
	flagUsage := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)

	var flagValue bool
	flag.BoolVar(
		&flagValue,
		flagName,
		isDefaultEnabled,
		flagUsage)
	collectorState[collector] = &flagValue

	factories[collector] = factory
}

// Exporter represents exporter and has a link to the client
type Exporter struct {
	collectors map[string]Collector
	client     *ftl_client.FTLClient
}

// NewExporter creates exporter using the provided socket path
func NewExporter(socket string) (*Exporter, error) {
	log.Printf("Initialize exporter using socket path: %s", socket)

	collectors := make(map[string]Collector)
	for key, enabled := range collectorState {
		if *enabled {
			collector, err := factories[key]()
			if err != nil {
				return nil, err
			}

			log.Println("Collector", key, "is enabled")

			collectors[key] = collector
		}
	}

	client, err := ftl_client.NewClient(socket)
	if err != nil {
		return nil, err
	}

	return &Exporter{
		collectors: collectors,
		client:     client,
	}, nil
}

// Describe implements the prometheus.Collector interface.
func (collector Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect implements the prometheus.Collector interface.
func (collector Exporter) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(collector.collectors))
	for name, c := range collector.collectors {
		go func(name string, c Collector) {
			execute(name, c, collector.client, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func execute(name string, c Collector, client *ftl_client.FTLClient, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.update(client, ch)
	duration := time.Since(begin)

	success := float64(1)
	if err != nil {
		success = 0
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	update(client *ftl_client.FTLClient, ch chan<- prometheus.Metric) error
}
