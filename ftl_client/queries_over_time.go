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

// GetQueriesOverTime retrieves amount of allowed and blocked queries
// for the last 24 hours aggregated over 10 minute intervals
// from response of `>overTime` command
func (client *FTLClient) GetQueriesOverTime() (*OverTime, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if err := sendCommand(conn, ">overTime"); err != nil {
		return nil, err
	}

	var lines struct {
		_     uint8
		Lines uint16
	}
	if err := binary.Read(conn, binary.BigEndian, &lines); err != nil {
		return nil, err
	}

	forwarded := make([]TimestampCount, lines.Lines)

	if err := binary.Read(conn, binary.BigEndian, &forwarded); err != nil {
		return nil, err
	}

	if err := binary.Read(conn, binary.BigEndian, &lines); err != nil {
		return nil, err
	}

	blocked := make([]TimestampCount, lines.Lines)

	if err := binary.Read(conn, binary.BigEndian, &blocked); err != nil {
		return nil, err
	}

	return &OverTime{
		Forwarded: forwarded,
		Blocked:   blocked,
	}, nil
}
