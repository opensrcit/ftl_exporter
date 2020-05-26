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
	"log"
	"net"
	"sort"
)

const (
	namespace = "ftl"
)

var (
	// >stats
	domainsBeingBlocked = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "domains_being_blocked"),
		"Domains being blocked.",
		nil, nil,
	)

	dnsQueriesToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dns_queries_today"),
		"DNS Queries today.",
		nil, nil,
	)

	adsBlockedToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ads_blocked_today"),
		"Ads blocked today.",
		nil, nil,
	)

	adsPercentageToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ads_percentage_today"),
		"Ads percentage today.",
		nil, nil,
	)

	uniqueDomainsToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "unique_domains_today"),
		"Unique domains seen today.",
		nil, nil,
	)

	queriesForwardedToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "queries_forwarded_today"),
		"Queries forwarded today.",
		nil, nil,
	)

	queriesCachedToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "queries_cached_today"),
		"Queries cached today.",
		nil, nil,
	)

	clientsEverSeen = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "clients_ever_seen"),
		"Clients ever seen.",
		nil, nil,
	)

	uniqueClients = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "unique_clients"),
		"Unique clients.",
		nil, nil,
	)

	status = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "status"),
		"Blocking status.",
		nil, nil,
	)

	// >top-domains
	overallQueriesToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "overall_queries_today"),
		"Overall queries today.",
		nil, nil,
	)

	topQueriesToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_queries_today"),
		"Top queries today.",
		[]string{"domain"}, nil,
	)

	// >top-ads
	topAdsToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_ads_today"),
		"Top Ads today.",
		[]string{"domain"}, nil,
	)

	overallAdsToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "overall_ads_today"),
		"Overall ads.",
		nil, nil,
	)

	// >top-clients
	topSourcesToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_sources_today"),
		"Top sources today.",
		[]string{"client"}, nil,
	)

	topBlockedSourcesToday = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "top_blocked_sources_today"),
		"Top blocked sources today.",
		[]string{"client"}, nil,
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

// Exporter represents exporter and has a link to the client
type Exporter struct {
	client *ftl_client.Client
}

// NewExporter creates exporter using the provided socket path
func NewExporter(socket string) *Exporter {
	log.Printf("Initialize exporter using socket path: %s", socket)

	c, err := net.Dial("unix", socket)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := c.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	client := ftl_client.NewClient(socket)

	return &Exporter{
		client: client,
	}
}

// Describe describes all the metrics ever exported by the exporter.
// It implements prometheus.Collector.
func (collector *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- domainsBeingBlocked
	ch <- dnsQueriesToday
	ch <- adsBlockedToday
	ch <- adsPercentageToday
	ch <- uniqueDomainsToday
	ch <- queriesForwardedToday
	ch <- queriesCachedToday
	ch <- clientsEverSeen
	ch <- uniqueClients
	ch <- status

	ch <- overallQueriesToday
	ch <- topQueriesToday

	ch <- topAdsToday
	ch <- overallAdsToday

	ch <- topSourcesToday
	ch <- topBlockedSourcesToday

	ch <- forwardDestinationsToday

	ch <- queryTypes

	ch <- queriesInDatabase
	ch <- databaseFilesize

	ch <- forwardedOverTime
	ch <- blockedOverTime

	ch <- clientsOverTimeMetric
}

// Collect is called by the Prometheus registry when collecting
// metrics.
func (collector *Exporter) Collect(ch chan<- prometheus.Metric) {
	stats, err := collector.client.GetStats()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	ch <- prometheus.MustNewConstMetric(domainsBeingBlocked, prometheus.GaugeValue, float64(stats.DomainsBeingBlocked.Value))
	ch <- prometheus.MustNewConstMetric(dnsQueriesToday, prometheus.GaugeValue, float64(stats.DnsQueries.Value))
	ch <- prometheus.MustNewConstMetric(adsBlockedToday, prometheus.GaugeValue, float64(stats.AdsBlocked.Value))
	ch <- prometheus.MustNewConstMetric(adsPercentageToday, prometheus.GaugeValue, float64(stats.AdsPercentage.Value))
	ch <- prometheus.MustNewConstMetric(uniqueDomainsToday, prometheus.GaugeValue, float64(stats.UniqueDomains.Value))
	ch <- prometheus.MustNewConstMetric(queriesForwardedToday, prometheus.GaugeValue, float64(stats.QueriesForwarded.Value))
	ch <- prometheus.MustNewConstMetric(queriesCachedToday, prometheus.GaugeValue, float64(stats.QueriesCached.Value))
	ch <- prometheus.MustNewConstMetric(clientsEverSeen, prometheus.GaugeValue, float64(stats.ClientsEverSeen.Value))
	ch <- prometheus.MustNewConstMetric(uniqueClients, prometheus.GaugeValue, float64(stats.UniqueClients.Value))
	ch <- prometheus.MustNewConstMetric(status, prometheus.GaugeValue, float64(stats.Status.Value))

	queries, err := collector.client.GetTopDomains()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	ch <- prometheus.MustNewConstMetric(overallQueriesToday, prometheus.GaugeValue, float64(queries.Total.Value))

	for _, hits := range queries.List {
		ch <- prometheus.MustNewConstMetric(topQueriesToday, prometheus.GaugeValue, float64(hits.Count.Value), hits.Domain)
	}

	ads, err := collector.client.GetTopAds()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	ch <- prometheus.MustNewConstMetric(overallAdsToday, prometheus.GaugeValue, float64(ads.Total.Value))

	for _, hits := range ads.List {
		ch <- prometheus.MustNewConstMetric(topAdsToday, prometheus.GaugeValue, float64(hits.Count.Value), hits.Domain)
	}

	clients, err := collector.client.GetTopClients()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	for _, hits := range clients.List {
		ch <- prometheus.MustNewConstMetric(topSourcesToday, prometheus.GaugeValue, float64(hits.Count.Value), hits.Domain)
	}

	blockedClients, err := collector.client.GetTopBlockedClients()
	if err != nil {
		log.Fatalf("failed to get data: %v", err)
	}

	for _, hits := range blockedClients.List {
		ch <- prometheus.MustNewConstMetric(topBlockedSourcesToday, prometheus.GaugeValue, float64(hits.Count.Value), hits.Domain)
	}

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
