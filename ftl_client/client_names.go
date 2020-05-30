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
	"io"
	"net"
)

// GetClientNames retrieves ordered list of client's names from
// response of `>client-names` command
func (client *FTLClient) GetClientNames() (*Clients, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(">client-names")); err != nil {
		return nil, err
	}

	var clients Clients
	for {
		var format uint8
		err := binary.Read(conn, binary.BigEndian, &format)

		if err == io.EOF || format == formatEOF {
			break
		}

		if err != nil {
			return nil, err
		}

		var length uint32

		err = binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			return nil, err
		}

		name := make([]byte, length)

		err = binary.Read(conn, binary.BigEndian, &name)
		if err != nil {
			return nil, err
		}

		err = binary.Read(conn, binary.BigEndian, &format)
		if err != nil {
			return nil, err
		}

		err = binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			return nil, err
		}

		address := make([]byte, length)

		err = binary.Read(conn, binary.BigEndian, &address)
		if err != nil {
			return nil, err
		}

		clients.List = append(clients.List, struct {
			Name    string
			Address string
		}{Name: string(name), Address: string(address)})
	}

	return &clients, nil
}
