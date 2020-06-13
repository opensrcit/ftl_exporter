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

// GetForwardDestinations retrieves forward destination with amount
// of queries forwarded to them from response of `>forward-dest` command
func (client *FTLClient) GetForwardDestinations() (*[]UpstreamDestination, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, ">forward-dest"); err != nil {
		return nil, err
	}

	var destinations []UpstreamDestination
	for {
		name, err := readString(conn)
		if err == errEndOfInput {
			break
		}
		if err != nil {
			return nil, err
		}

		address, err := readString(conn)
		if err != nil {
			return nil, err
		}

		percentage, err := readFloat32(conn)
		if err != nil {
			return nil, err
		}

		destinations = append(destinations, UpstreamDestination{
			Name:       name,
			Address:    address,
			Percentage: percentage,
		})
	}

	return &destinations, nil
}
