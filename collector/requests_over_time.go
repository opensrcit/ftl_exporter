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
	"sort"
)

type requestsOverTimeCollector struct {
	requestsForwarded *prometheus.Desc
	requestsBlocked   *prometheus.Desc
}

func init() {
	registerCollector("requests_over_time", defaultEnabled, newRequestsOverTimeCollector)
}

func newRequestsOverTimeCollector() (Collector, error) {
	return &requestsOverTimeCollector{
		requestsForwarded: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "requests_allowed"),
			"Amount of allowed requests for the last 10 minutes.",
			nil, nil,
		),

		requestsBlocked: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "requests_blocked"),
			"Amount blocked requests for the last 10 minutes.",
			nil, nil,
		),
	}, nil
}

func (c *requestsOverTimeCollector) update(client *ftl_client.Client, ch chan<- prometheus.Metric) error {
	queriesOverTime, err := client.GetQueriesOverTime()
	if err != nil {
		return err
	}

	sort.SliceStable(queriesOverTime.Forwarded, func(i, j int) bool {
		return queriesOverTime.Forwarded[i].Timestamp.Value > queriesOverTime.Forwarded[j].Timestamp.Value
	})
	lastForwardedOverTime := queriesOverTime.Forwarded[:1]
	for _, hits := range lastForwardedOverTime {
		ch <- prometheus.MustNewConstMetric(c.requestsForwarded, prometheus.GaugeValue, float64(hits.Count.Value))
	}

	sort.SliceStable(queriesOverTime.Blocked, func(i, j int) bool {
		return queriesOverTime.Blocked[i].Timestamp.Value > queriesOverTime.Blocked[j].Timestamp.Value
	})
	lastBlockedOverTime := queriesOverTime.Blocked[:1]
	for _, hits := range lastBlockedOverTime {
		ch <- prometheus.MustNewConstMetric(c.requestsBlocked, prometheus.GaugeValue, float64(hits.Count.Value))
	}

	return nil
}
