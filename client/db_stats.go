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

// GetDBStats retrieves database statistics from response of `>dbstats` command
func (client *FTLClient) GetDBStats() (*DBStats, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, ">dbstats"); err != nil {
		return nil, err
	}

	var stats struct {
		Rows ftlInt32
		Size ftlUInt64
	}
	if err := binary.Read(conn, binary.BigEndian, &stats); err != nil {
		return nil, err
	}

	return &DBStats{
		RowsCount: int(stats.Rows.Value),
		FileSize:  int(stats.Size.Value),
	}, nil
}
