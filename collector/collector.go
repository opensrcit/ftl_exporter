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
	"sort"
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

	// >forward-dest
	forwardDestinationsToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "forward_destinations_today"),
		"Forward destinations today.",
		[]string{"address"}, nil,
	)

	// >querytypes
	queryTypes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "query_types_today"),
		"DNS Query types today.",
		[]string{"query"}, nil,
	)

	// >dbstats
	queriesInDatabase = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "queries_in_database"),
		"Queries in database.",
		nil, nil,
	)

	databaseFilesize = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "database_filesize"),
		"Database file size.",
		nil, nil,
	)

	// >overTime
	forwardedOverTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "forwarded_over_time"),
		"Forwarded queries over time (last 10 minutes).",
		nil, nil,
	)

	blockedOverTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "blocked_over_time"),
		"Blocked queries over time (last 10 minutes).",
		nil, nil,
	)

	// >ClientsoverTime TODO: is it public api?
	clientsOverTimeMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "clients_over_time"),
		"Client requests over time (last 10 minutes).",
		[]string{"address"}, nil,
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
	client     *ftl_client.Client
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

// Collect is called by the Prometheus registry when collecting
// metrics.
func (collector *Exporter) Collect2(ch chan<- prometheus.Metric) {
	destinations, err := collector.client.GetForwardDestinations()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	for _, hits := range *destinations {
		ch <- prometheus.MustNewConstMetric(forwardDestinationsToday, prometheus.GaugeValue, float64(hits.Percentage.Value), hits.Address)
	}

	queryTypesData, err := collector.client.GetQueryTypes()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	for _, hits := range *queryTypesData {
		ch <- prometheus.MustNewConstMetric(queryTypes, prometheus.GaugeValue, float64(hits.Percentage.Value), hits.Entry)
	}

	dbStats, err := collector.client.GetDBStats()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	ch <- prometheus.MustNewConstMetric(queriesInDatabase, prometheus.CounterValue, float64(dbStats.Rows.Value))
	ch <- prometheus.MustNewConstMetric(databaseFilesize, prometheus.CounterValue, float64(dbStats.Size.Value))

	queriesOverTime, err := collector.client.GetQueriesOverTime()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	sort.SliceStable(queriesOverTime.Forwarded, func(i, j int) bool {
		return queriesOverTime.Forwarded[i].Timestamp.Value > queriesOverTime.Forwarded[j].Timestamp.Value
	})
	lastForwardedOverTime := queriesOverTime.Forwarded[:1]
	for _, hits := range lastForwardedOverTime {
		ch <- prometheus.MustNewConstMetric(forwardedOverTime, prometheus.GaugeValue, float64(hits.Count.Value))
	}

	sort.SliceStable(queriesOverTime.Blocked, func(i, j int) bool {
		return queriesOverTime.Blocked[i].Timestamp.Value > queriesOverTime.Blocked[j].Timestamp.Value
	})
	lastBlockedOverTime := queriesOverTime.Blocked[:1]
	for _, hits := range lastBlockedOverTime {
		ch <- prometheus.MustNewConstMetric(blockedOverTime, prometheus.GaugeValue, float64(hits.Count.Value))
	}

	clientsOverTime, err := collector.client.GetClientsOverTime()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	clientNames, err := collector.client.GetClientNames()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	sort.SliceStable(clientsOverTime.List, func(i, j int) bool {
		return clientsOverTime.List[i].Timestamp.Value > clientsOverTime.List[j].Timestamp.Value
	})
	lastClientsOverTime := clientsOverTime.List[:1]
	for _, hits := range lastClientsOverTime {
		for i, count := range hits.Count {
			address := fmt.Sprintf("address_%d", i)
			if i < len(clientNames.List) {
				address = clientNames.List[i].Address
			}
			ch <- prometheus.MustNewConstMetric(
				clientsOverTimeMetric,
				prometheus.GaugeValue,
				float64(count.Value),
				address,
			)
		}
	}
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

func execute(name string, c Collector, client *ftl_client.Client, ch chan<- prometheus.Metric) {
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
	update(client *ftl_client.Client, ch chan<- prometheus.Metric) error
}
