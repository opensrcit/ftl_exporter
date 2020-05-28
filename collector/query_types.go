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

type queryTypesCollector struct {
	queryTypesToday *prometheus.Desc
}

func init() {
	registerCollector("query_types", defaultEnabled, newQueryTypesCollector)
}

// newDomainCollector returns a new Collector exposing >querytypes command
func newQueryTypesCollector() (Collector, error) {
	return &queryTypesCollector{
		queryTypesToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "query_types_today"),
			"DNS Query types today (percentage).",
			[]string{"query"}, nil,
		),
	}, nil
}

// update implements Collector and exposes metrics from >querytypes command
func (c *queryTypesCollector) update(client *ftl_client.Client, ch chan<- prometheus.Metric) error {
	queryTypesData, err := client.GetQueryTypes()
	if err != nil {
		return err
	}

	for _, hits := range *queryTypesData {
		ch <- prometheus.MustNewConstMetric(c.queryTypesToday, prometheus.GaugeValue, float64(hits.Percentage.Value), hits.Entry)
	}

	return nil
}
