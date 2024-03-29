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

type forwardDestinationCollector struct {
	forwardDestinationsToday *prometheus.Desc
}

func init() {
	registerCollector("forward_destinations", defaultEnabled, newForwardDestinationCollector)
}

func newForwardDestinationCollector() (Collector, error) {
	return &forwardDestinationCollector{
		forwardDestinationsToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "forward_destinations_today"),
			"Forward destinations today.",
			[]string{"address"}, nil,
		),
	}, nil
}

func (c *forwardDestinationCollector) update(client *client.FTLClient, ch chan<- prometheus.Metric) error {
	destinations, err := client.GetForwardDestinations()
	if err != nil {
		return err
	}

	for _, hits := range *destinations {
		ch <- prometheus.MustNewConstMetric(c.forwardDestinationsToday, prometheus.GaugeValue, float64(hits.Percentage), hits.Address)
	}

	return nil
}
