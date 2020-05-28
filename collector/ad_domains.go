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

type adDomainCollector struct {
	totalAdDomainsToday *prometheus.Desc
	topAdDomainsToday   *prometheus.Desc
}

func init() {
	registerCollector("ad_domains", defaultEnabled, newAdDomainCollector)
}

// newDomainCollector returns a new Collector exposing >top-ads command
func newAdDomainCollector() (Collector, error) {
	return &adDomainCollector{
		totalAdDomainsToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "total_ad_domains_today"),
			"Overall ads.",
			nil, nil,
		),

		topAdDomainsToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "top_ad_domains_today"),
			"Top Ads today.",
			[]string{"domain"}, nil,
		),
	}, nil
}

// update implements Collector and exposes metrics from >top-ads command
func (c *adDomainCollector) update(client *ftl_client.Client, ch chan<- prometheus.Metric) error {
	queries, err := client.GetTopAds()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.totalAdDomainsToday, prometheus.GaugeValue, float64(queries.Total.Value))

	for _, hits := range queries.List {
		ch <- prometheus.MustNewConstMetric(c.topAdDomainsToday, prometheus.GaugeValue, float64(hits.Count.Value), hits.Domain)
	}

	return nil
}
