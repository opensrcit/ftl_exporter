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
	"fmt"
	"github.com/opensrcit/ftl_exporter/ftl_client"
	"github.com/prometheus/client_golang/prometheus"
	"sort"
)

type clientsOverTimeCollector struct {
	clients *prometheus.Desc
}

func init() {
	// command >ClientsoverTime is not in the official api
	// it is disabled by default
	registerCollector("clients_over_time", defaultDisabled, newClientsOverTimeCollector)
}

func newClientsOverTimeCollector() (Collector, error) {
	return &clientsOverTimeCollector{
		clients: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "clients"),
			"Client requests for the last 10 minutes.",
			[]string{"address"}, nil,
		),
	}, nil
}

func (c *clientsOverTimeCollector) update(client *ftl_client.FTLClient, ch chan<- prometheus.Metric) error {
	clientsOverTime, err := client.GetClientsOverTime()
	if err != nil {
		return err
	}

	clientNames, err := client.GetClientNames()
	if err != nil {
		return err
	}

	sort.SliceStable(*clientsOverTime, func(i, j int) bool {
		return (*clientsOverTime)[i].Timestamp > (*clientsOverTime)[j].Timestamp
	})
	lastClientsOverTime := (*clientsOverTime)[:1]
	for _, hits := range lastClientsOverTime {
		for i, count := range hits.Count {
			address := fmt.Sprintf("address_%d", i)
			if i < len(*clientNames) {
				address = (*clientNames)[i].Address
			}
			ch <- prometheus.MustNewConstMetric(
				c.clients,
				prometheus.GaugeValue,
				float64(count.Value),
				address,
			)
		}
	}

	return nil
}
