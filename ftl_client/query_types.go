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

// GetQueryTypes retrieves query type and their percentages among all queries
// from response of `>querytypes` command
func (client *FTLClient) GetQueryTypes() (*[]PercentageEntry, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(">querytypes")); err != nil {
		return nil, err
	}

	var queryTypes []PercentageEntry
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

		var percentage Float32Block

		err = binary.Read(conn, binary.BigEndian, &percentage)
		if err != nil {
			return nil, err
		}

		queryTypes = append(queryTypes, PercentageEntry{
			Entry:      string(name),
			Percentage: percentage,
		})
	}

	return &queryTypes, nil
}
