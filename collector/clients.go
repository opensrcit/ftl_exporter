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
	"github.com/opensrcit/ftl_exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type clientCollector struct {
	topClientsToday        *prometheus.Desc
	topBlockedClientsToday *prometheus.Desc
}

func init() {
	registerCollector("clients", defaultEnabled, newClientCollector)
}

func newClientCollector() (Collector, error) {
	return &clientCollector{
		topClientsToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "top_clients_today"),
			"Top sources today.",
			[]string{"client"}, nil,
		),

		topBlockedClientsToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "top_blocked_clients_today"),
			"Top blocked sources today.",
			[]string{"client"}, nil,
		),
	}, nil
}

func (c *clientCollector) update(client *client.FTLClient, ch chan<- prometheus.Metric) error {
	clients, err := client.GetTopClients()
	if err != nil {
		return err
	}

	for _, hits := range clients.List {
		ch <- prometheus.MustNewConstMetric(c.topClientsToday, prometheus.GaugeValue, float64(hits.Count), hits.Entry)
	}

	blockedClients, err := client.GetTopBlockedClients()
	if err != nil {
		return err
	}

	for _, hits := range blockedClients.List {
		ch <- prometheus.MustNewConstMetric(c.topBlockedClientsToday, prometheus.GaugeValue, float64(hits.Count), hits.Entry)
	}

	return nil
}
