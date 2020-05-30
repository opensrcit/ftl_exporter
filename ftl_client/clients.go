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

package ftl_client

import (
	"encoding/binary"
	"net"
)

// GetTopClients retrieves the list of clients together with amount of queries
// made by each client from response of `>top-clients` command
func (client *FTLClient) GetTopClients() (*Entries, error) {
	return topClientsFor(">top-clients", client)
}

// GetTopBlockedClients retrieves the list of clients together with amount of blocked
// queries made by each client from response of `>top-clients` command
func (client *FTLClient) GetTopBlockedClients() (*Entries, error) {
	return topClientsFor(">top-clients blocked", client)
}

func topClientsFor(command string, client *FTLClient) (*Entries, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(command)); err != nil {
		return nil, err
	}

	var result Entries

	if err := binary.Read(conn, binary.BigEndian, &result.Total); err != nil {
		return nil, err
	}

	for {
		_, err := readString(conn)
		if err == EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		address, err := readString(conn)
		if err != nil {
			return nil, err
		}

		count, err := readUint32(conn)
		if err != nil {
			return nil, err
		}

		result.List = append(result.List, struct {
			Entry string
			Count uint32
		}{Entry: address, Count: count})
	}

	return &result, nil
}
