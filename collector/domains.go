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
	"github.com/opensrcit/ftl_exporter/ftl_client"
	"github.com/prometheus/client_golang/prometheus"
)

type domainCollector struct {
	totalDomainsToday *prometheus.Desc
	topDomainsToday   *prometheus.Desc
}

func init() {
	registerCollector("domains", defaultEnabled, newDomainCollector)
}

func newDomainCollector() (Collector, error) {
	return &domainCollector{
		totalDomainsToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "total_domains_today"),
			"Total domains today.",
			nil, nil,
		),

		topDomainsToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "top_domains_today"),
			"Top domains today.",
			[]string{"domain"}, nil,
		),
	}, nil
}

func (c *domainCollector) update(client *ftl_client.Client, ch chan<- prometheus.Metric) error {
	queries, err := client.GetTopDomains()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.totalDomainsToday, prometheus.GaugeValue, float64(queries.Total.Value))

	for _, hits := range queries.List {
		ch <- prometheus.MustNewConstMetric(c.topDomainsToday, prometheus.GaugeValue, float64(hits.Count.Value), hits.Domain)
	}

	return nil
}
