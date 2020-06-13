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

// GetTopDomains retrieves the list of domains together with amount of queries
// made for each domain from response of `>top-domains` command
func (client *FTLClient) GetTopDomains() (*TopEntries, error) {
	return topQueriesFor(">top-domains", client)
}

// GetTopAds retrieves the list of ad domains together with amount of queries
// made for each domain from response of `>top-ads` command
func (client *FTLClient) GetTopAds() (*TopEntries, error) {
	return topQueriesFor(">top-ads", client)
}

func topQueriesFor(command string, client *FTLClient) (*TopEntries, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, command); err != nil {
		return nil, err
	}

	var total ftlInt32
	if err := binary.Read(conn, binary.BigEndian, &total); err != nil {
		return nil, err
	}

	result := TopEntries{
		Total: int(total.Value),
	}

	for {
		domainName, err := readString(conn)
		if err == errEndOfInput {
			break
		}
		if err != nil {
			return nil, err
		}

		domainCount, err := readInt32(conn)
		if err != nil {
			return nil, err
		}

		result.Entries = append(result.Entries, struct {
			Label string
			Count int
		}{Label: domainName, Count: int(domainCount)})
	}

	return &result, nil
}
