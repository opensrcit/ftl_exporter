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
	"net"
)

// GetQueryTypes retrieves map with query type as keys and their percentages
// among all queries as values from response of `>querytypes` command
func (client *FTLClient) GetQueryTypes() (*map[string]float32, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, ">querytypes"); err != nil {
		return nil, err
	}

	queryTypes := make(map[string]float32)
	for {
		name, err := readString(conn)
		if err == errEndOfInput {
			break
		}
		if err != nil {
			return nil, err
		}

		percentage, err := readFloat32(conn)
		if err != nil {
			return nil, err
		}

		queryTypes[name] = percentage
	}

	return &queryTypes, nil
}
