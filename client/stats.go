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

import (
	"encoding/binary"
	"net"
)

// GetStats retrieves engine statistics from response of `>stats` command
func (client *FTLClient) GetStats() (*Stats, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, ">stats"); err != nil {
		return nil, err
	}

	var stats struct {
		DomainsBeingBlocked UInt32Block
		DnsQueriesToday     UInt32Block
		AdsBlockedToday     UInt32Block
		AdsPercentageToday  Float32Block
		UniqueDomains       UInt32Block
		QueriesForwarded    UInt32Block
		QueriesCached       UInt32Block
		ClientsEverSeen     UInt32Block
		UniqueClients       UInt32Block
		Status              UInt8Block
	}
	if err := binary.Read(conn, binary.BigEndian, &stats); err != nil {
		return nil, err
	}

	return &Stats{
		DomainsBeingBlocked: int(stats.DomainsBeingBlocked.Value),
		DnsQueriesToday:     int(stats.DnsQueriesToday.Value),
		AdsBlockedToday:     int(stats.AdsBlockedToday.Value),
		AdsPercentageToday:  stats.AdsPercentageToday.Value,
		UniqueDomains:       int(stats.UniqueDomains.Value),
		QueriesForwarded:    int(stats.QueriesForwarded.Value),
		QueriesCached:       int(stats.QueriesCached.Value),
		ClientsEverSeen:     int(stats.ClientsEverSeen.Value),
		UniqueClients:       int(stats.UniqueClients.Value),
		Status:              int(stats.Status.Value),
	}, nil
}
