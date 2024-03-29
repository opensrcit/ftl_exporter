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

// DBStats represents the response of `>db-stats` command
type DBStats struct {
	Rows UInt32Block
	Size UInt64Block
}

type UInt32Block struct {
	_     uint8
	Value uint32
}

type Int32Block struct {
	_     uint8
	Value int32
}

type UInt64Block struct {
	_     uint8
	Value uint64
}

type UInt8Block struct {
	_     uint8
	Value uint8
}

type Float32Block struct {
	_     uint8
	Value float32
}

type Entries struct {
	Total UInt32Block
	List  []struct {
		Entry string
		Count uint32
	}
}

type UpstreamDestination struct {
	Name       string
	Address    string
	Percentage float32
}

type TimestampCount struct {
	Timestamp UInt32Block
	Count     UInt32Block
}

type TimestampClients struct {
	Timestamp uint32
	Count     []Int32Block
}

type ClientsOverTime struct {
	List []TimestampClients
}

type Client struct {
	Name    string
	Address string
}

type OverTime struct {
	Forwarded []TimestampCount
	Blocked   []TimestampCount
}
