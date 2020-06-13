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
	"sort"
)

type queriesOverTimeCollector struct {
	queriesForwarded *prometheus.Desc
	queriesBlocked   *prometheus.Desc
}

func init() {
	registerCollector("queries_over_time", defaultEnabled, newQueriesOverTimeCollector)
}

func newQueriesOverTimeCollector() (Collector, error) {
	return &queriesOverTimeCollector{
		queriesForwarded: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "queries_allowed"),
			"Amount of allowed queries for the last 10 minutes.",
			nil, nil,
		),

		queriesBlocked: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "queries_blocked"),
			"Amount blocked queries for the last 10 minutes.",
			nil, nil,
		),
	}, nil
}

func (c *queriesOverTimeCollector) update(client *client.FTLClient, ch chan<- prometheus.Metric) error {
	queriesOverTime, err := client.GetQueriesOverTime()
	if err != nil {
		return err
	}

	sort.SliceStable(queriesOverTime.Forwarded, func(i, j int) bool {
		return queriesOverTime.Forwarded[i].Timestamp > queriesOverTime.Forwarded[j].Timestamp
	})
	lastForwardedOverTime := queriesOverTime.Forwarded[:1]
	for _, hits := range lastForwardedOverTime {
		ch <- prometheus.MustNewConstMetric(c.queriesForwarded, prometheus.GaugeValue, float64(hits.Count))
	}

	sort.SliceStable(queriesOverTime.Blocked, func(i, j int) bool {
		return queriesOverTime.Blocked[i].Timestamp > queriesOverTime.Blocked[j].Timestamp
	})
	lastBlockedOverTime := queriesOverTime.Blocked[:1]
	for _, hits := range lastBlockedOverTime {
		ch <- prometheus.MustNewConstMetric(c.queriesBlocked, prometheus.GaugeValue, float64(hits.Count))
	}

	return nil
}
