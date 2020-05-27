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

type statsCollector struct {
	domainsBeingBlocked   *prometheus.Desc
	dnsQueriesToday       *prometheus.Desc
	adsBlockedToday       *prometheus.Desc
	adsPercentageToday    *prometheus.Desc
	uniqueDomainsToday    *prometheus.Desc
	queriesForwardedToday *prometheus.Desc
	queriesCachedToday    *prometheus.Desc
	clientsEverSeen       *prometheus.Desc
	uniqueClients         *prometheus.Desc
	status                *prometheus.Desc
}

func init() {
	registerCollector("stats", defaultEnabled, newStatsCollector)
}

// newStatsCollector returns a new Collector exposing >stats
func newStatsCollector() (Collector, error) {
	return &statsCollector{
		domainsBeingBlocked: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "domains_being_blocked"),
			"Domains being blocked.",
			nil, nil,
		),

		dnsQueriesToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "dns_queries_today"),
			"DNS Queries today.",
			nil, nil,
		),
		adsBlockedToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "ads_blocked_today"),
			"Ads blocked today.",
			nil, nil,
		),
		adsPercentageToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "ads_percentage_today"),
			"Ads percentage today.",
			nil, nil,
		),
		uniqueDomainsToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unique_domains_today"),
			"Unique domains seen today.",
			nil, nil,
		),
		queriesForwardedToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "queries_forwarded_today"),
			"Queries forwarded today.",
			nil, nil,
		),
		queriesCachedToday: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "queries_cached_today"),
			"Queries cached today.",
			nil, nil,
		),
		clientsEverSeen: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "clients_ever_seen"),
			"Clients ever seen.",
			nil, nil,
		),
		uniqueClients: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "unique_clients"),
			"Unique clients.",
			nil, nil,
		),
		status: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "status"),
			"Blocking status.",
			nil, nil,
		),
	}, nil
}

// update implements Collector and exposes metrics from >stats command
func (c *statsCollector) update(client *ftl_client.Client, ch chan<- prometheus.Metric) error {
	stats, err := client.GetStats()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.domainsBeingBlocked, prometheus.GaugeValue, float64(stats.DomainsBeingBlocked.Value))
	ch <- prometheus.MustNewConstMetric(c.dnsQueriesToday, prometheus.GaugeValue, float64(stats.DnsQueries.Value))
	ch <- prometheus.MustNewConstMetric(c.adsBlockedToday, prometheus.GaugeValue, float64(stats.AdsBlocked.Value))
	ch <- prometheus.MustNewConstMetric(c.adsPercentageToday, prometheus.GaugeValue, float64(stats.AdsPercentage.Value))
	ch <- prometheus.MustNewConstMetric(c.uniqueDomainsToday, prometheus.GaugeValue, float64(stats.UniqueDomains.Value))
	ch <- prometheus.MustNewConstMetric(c.queriesForwardedToday, prometheus.GaugeValue, float64(stats.QueriesForwarded.Value))
	ch <- prometheus.MustNewConstMetric(c.queriesCachedToday, prometheus.GaugeValue, float64(stats.QueriesCached.Value))
	ch <- prometheus.MustNewConstMetric(c.clientsEverSeen, prometheus.GaugeValue, float64(stats.ClientsEverSeen.Value))
	ch <- prometheus.MustNewConstMetric(c.uniqueClients, prometheus.GaugeValue, float64(stats.UniqueClients.Value))
	ch <- prometheus.MustNewConstMetric(c.status, prometheus.GaugeValue, float64(stats.Status.Value))

	return nil
}
