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
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	addr *net.UnixAddr
}

const formatEOF uint8 = 0xc1 // 193

type Stats struct {
	DomainsBeingBlocked UInt32Block
	DnsQueries          UInt32Block
	AdsBlocked          UInt32Block
	AdsPercentage       Float32Block
	UniqueDomains       UInt32Block
	QueriesForwarded    UInt32Block
	QueriesCached       UInt32Block
	ClientsEverSeen     UInt32Block
	UniqueClients       UInt32Block
	Status              UInt8Block
}

type DBStats struct {
	Rows UInt32Block
	Size UInt64Block
}

type UInt32Block struct {
	_     uint8
	Value uint32
}

type Int32Block struct {
	_     uint8
	Value int32
}

type UInt64Block struct {
	_     uint8
	Value uint64
}

type UInt8Block struct {
	_     uint8
	Value uint8
}

type Float32Block struct {
	_     uint8
	Value float32
}

type DomainEntries struct {
	Total UInt32Block
	List  []struct {
		Domain string
		Count  UInt32Block
	}
}

type PercentageEntry struct {
	Entry      string
	Percentage Float32Block
}

type UpstreamDestination struct {
	Name       string
	Address    string
	Percentage Float32Block
}

type TimestampCount struct {
	Timestamp UInt32Block
	Count     UInt32Block
}

type TimestampClients struct {
	Timestamp UInt32Block
	Count     []Int32Block
}

type ClientsOverTime struct {
	List []TimestampClients
}

type Clients struct {
	List []struct {
		Name    string
		Address string
	}
}

type OverTime struct {
	Forwarded []TimestampCount
	Blocked   []TimestampCount
}

func NewClient(socket string) *Client {
	addr, err := net.ResolveUnixAddr("unix", socket)
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}

	return &Client{
		addr: addr,
	}
}

func (client *Client) GetStats() (*Stats, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(">stats")); err != nil {
		return nil, err
	}

	var stats Stats
	if err := binary.Read(conn, binary.BigEndian, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (client *Client) GetTopDomains() (*DomainEntries, error) {
	return topQueriesFor(">top-domains", client)
}

func (client *Client) GetTopAds() (*DomainEntries, error) {
	return topQueriesFor(">top-ads", client)
}

func (client *Client) GetTopClients() (*DomainEntries, error) {
	return topClientsFor(">top-clients", client)
}

func (client *Client) GetTopBlockedClients() (*DomainEntries, error) {
	return topClientsFor(">top-clients blocked", client)
}

func (client *Client) GetForwardDestinations() (*[]UpstreamDestination, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(">forward-dest")); err != nil {
		return nil, err
	}

	var destinations []UpstreamDestination
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

		var percentage Float32Block

		err = binary.Read(conn, binary.BigEndian, &percentage)
		if err != nil {
			return nil, err
		}

		destinations = append(destinations, UpstreamDestination{
			Name:       string(name),
			Address:    string(address),
			Percentage: percentage,
		})
	}

	return &destinations, nil
}

func (client *Client) GetQueryTypes() (*[]PercentageEntry, error) {
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

func topClientsFor(command string, client *Client) (*DomainEntries, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(command)); err != nil {
		return nil, err
	}

	var result DomainEntries

	if err := binary.Read(conn, binary.BigEndian, &result.Total); err != nil {
		return nil, err
	}

	for {
		var format uint8
		err := binary.Read(conn, binary.BigEndian, &format)

		if err == io.EOF || format == formatEOF {
			break
		}

		if err != nil {
			return nil, err
		}

		var emptyString struct {
			_ uint32
			_ [0]byte
		}

		if err := binary.Read(conn, binary.BigEndian, &emptyString); err != nil {
			return nil, err
		}

		err = binary.Read(conn, binary.BigEndian, &format)
		if err != nil {
			return nil, err
		}

		var length uint32

		err = binary.Read(conn, binary.BigEndian, &length)
		if err != nil {
			return nil, err
		}

		address := make([]byte, length)

		err = binary.Read(conn, binary.BigEndian, &address)
		if err != nil {
			return nil, err
		}

		var count UInt32Block

		err = binary.Read(conn, binary.BigEndian, &count)
		if err != nil {
			return nil, err
		}

		result.List = append(result.List, struct {
			Domain string
			Count  UInt32Block
		}{Domain: string(address), Count: count})
	}

	return &result, nil
}

func topQueriesFor(command string, client *Client) (*DomainEntries, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(command)); err != nil {
		return nil, err
	}

	var result DomainEntries
	if err := binary.Read(conn, binary.BigEndian, &result.Total); err != nil {
		return nil, err
	}

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

		domainName := make([]byte, length)

		err = binary.Read(conn, binary.BigEndian, &domainName)
		if err != nil {
			return nil, err
		}

		var domainCount UInt32Block

		err = binary.Read(conn, binary.BigEndian, &domainCount)
		if err != nil {
			return nil, err
		}

		result.List = append(result.List, struct {
			Domain string
			Count  UInt32Block
		}{Domain: string(domainName), Count: domainCount})
	}

	return &result, nil
}

func (client *Client) GetDBStats() (*DBStats, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(">dbstats")); err != nil {
		return nil, err
	}

	var stats DBStats
	if err := binary.Read(conn, binary.BigEndian, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (client *Client) GetQueriesOverTime() (*OverTime, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(">overTime")); err != nil {
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

func (client *Client) GetClientsOverTime() (*ClientsOverTime, error) {
	conn, err := net.DialUnix("unix", nil, client.addr)
	if err != nil {
		return nil, err
	}
	defer closeConnection(conn)

	if _, err := conn.Write([]byte(">ClientsoverTime")); err != nil {
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
			Timestamp: UInt32Block{
				Value: timestamp,
			},
			Count: clients,
		})
	}

	return &ClientsOverTime{
		List: timestamps,
	}, nil
}

func (client *Client) GetClientNames() (*Clients, error) {
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

func closeConnection(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
