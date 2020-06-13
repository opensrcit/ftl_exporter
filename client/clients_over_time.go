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
	"io"
	"net"
)

// GetClientsOverTime retrieves amount of queries grouped by client
// for the last 24 hours aggregated over 10 minute intervals
// from response of `>ClientsoverTime` command
// Warning: API might be not public
func (client *FTLClient) GetClientsOverTime() (*[]TimestampClients, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, ">ClientsoverTime"); err != nil {
		return nil, err
	}

	var timestamps []TimestampClients
	for {
		var format uint8
		err := binary.Read(conn, binary.BigEndian, &format)

		if err == io.EOF || format == formatEOF {
			break
		}

		var clients []Int32Block

		var timestamp uint32
		err = binary.Read(conn, binary.BigEndian, &timestamp)
		if err != nil {
			return nil, err
		}

		for {
			var clientQueryCount Int32Block
			err := binary.Read(conn, binary.BigEndian, &clientQueryCount)
			if err != nil {
				return nil, err
			}

			if clientQueryCount.Value == -1 {
				break
			}

			clients = append(clients, clientQueryCount)
		}

		timestamps = append(timestamps, TimestampClients{
			Timestamp: timestamp,
			Count:     clients,
		})
	}

	return &timestamps, nil
}
