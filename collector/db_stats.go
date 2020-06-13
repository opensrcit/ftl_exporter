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

type dbStatsCollector struct {
	queriesInDatabase *prometheus.Desc
	databaseFileSize  *prometheus.Desc
}

func init() {
	// >dbstats request may take some time for processing in case of a large database file
	// it is disabled by default
	registerCollector("db_stats", defaultDisabled, newDbStatsCollector)
}

func newDbStatsCollector() (Collector, error) {
	return &dbStatsCollector{
		queriesInDatabase: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "queries_in_database"),
			"Queries in database.",
			nil, nil,
		),

		databaseFileSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "database_file_size"),
			"Database file size.",
			nil, nil,
		),
	}, nil
}

func (c *dbStatsCollector) update(client *ftl_client.FTLClient, ch chan<- prometheus.Metric) error {
	dbStats, err := client.GetDBStats()
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.queriesInDatabase, prometheus.CounterValue, float64(dbStats.Rows.Value))
	ch <- prometheus.MustNewConstMetric(c.databaseFileSize, prometheus.CounterValue, float64(dbStats.Size.Value))

	return nil
}
