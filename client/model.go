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

package client

// Stats represents the response of `>stats` command
type Stats struct {
	DomainsBeingBlocked int
	DnsQueriesToday     int
	AdsBlockedToday     int
	AdsPercentageToday  float32
	UniqueDomains       int
	QueriesForwarded    int
	QueriesCached       int
	ClientsEverSeen     int
	UniqueClients       int
	Status              int
}

// DBStats represents the response of `>db-stats` command.
// It contains amount of rows in database and current file size of the database
type DBStats struct {
	RowsCount int
	FileSize  int
}

// TopEntries represents the response of `>top-clients` and `>top-domains` commands.
// It contains a total amount of entries and a list of entries label and count
type TopEntries struct {
	Total   int
	Entries []struct {
		Label string
		Count int
	}
}

// UpstreamDestination represents the response `>forward-dest` command.
// It contains a name, address and percentage of total requests
type UpstreamDestination struct {
	Name       string
	Address    string
	Percentage float32
}

// TimestampClients represents the response `>ClientsoverTime` command.
// It contains a timestamp and a list of amount of requests made by each client.
// Order of requests counts represents clients from `>client-names` command
type TimestampClients struct {
	Timestamp int
	Count     []int
}

// Client represents the response `>client-names` command.
// It contains a name and address of the client
type Client struct {
	Name    string
	Address string
}

// QueriesOverTime represents the response `>overTime` command.
// It contains list of amounts of forwarded and blocked requests grouped by 10 minute intervals
type QueriesOverTime struct {
	Forwarded []timestampCount
	Blocked   []timestampCount
}

type timestampCount struct {
	Timestamp int
	Count     int
}

type ftlUInt64 struct {
	_     uint8
	Value uint64
}

type ftlInt32 struct {
	_     uint8
	Value int32
}

type ftlInt8 struct {
	_     uint8
	Value uint8
}

type ftlFloat32 struct {
	_     uint8
	Value float32
}
