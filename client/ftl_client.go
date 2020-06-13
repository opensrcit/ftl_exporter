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
	"errors"
	"io"
	"log"
	"net"
)

const (
	formatInt32   uint8 = 0xd2 // 210
	formatFloat32 uint8 = 0xca // 202
	formatUInt8   uint8 = 0xcc // 204
	formatString  uint8 = 0xdb // 219
	formatMap16   uint8 = 0xde // 222

	formatEOF uint8 = 0xc1 // 193
)

var errEndOfInput = errors.New("end of the input")
var errInvalidFormat = errors.New("unexpected format")

// FTLClient for Pi-holes's FTL daemon. Contains address to a unix socket
type FTLClient struct {
	addr *net.UnixAddr
}

// NewClient creates the Pi-hole's FTL engine client
func NewClient(socket string) (*FTLClient, error) {
	addr, err := net.ResolveUnixAddr("unix", socket)
	if err != nil {
		return nil, err
	}

	c, err := net.Dial("unix", socket)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := c.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return &FTLClient{
		addr: addr,
	}, nil
}

func readString(conn *net.UnixConn) (string, error) {
	var format uint8
	if err := binary.Read(conn, binary.BigEndian, &format); err != nil {
		if err == io.EOF {
			return "", errEndOfInput
		}

		return "", err
	}

	if format == formatEOF {
		return "", errEndOfInput
	}

	if format != formatString {
		return "", errInvalidFormat
	}

	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return "", err
	}

	value := make([]byte, length)

	if err := binary.Read(conn, binary.BigEndian, &value); err != nil {
		return "", err
	}

	return string(value), nil
}

func readFloat32(conn *net.UnixConn) (float32, error) {
	var format uint8
	if err := binary.Read(conn, binary.BigEndian, &format); err != nil {
		if err == io.EOF {
			return 0.0, errEndOfInput
		}

		return 0.0, err
	}

	if format == formatEOF {
		return 0.0, errEndOfInput
	}

	if format != formatFloat32 {
		return 0.0, errInvalidFormat
	}

	var value float32
	if err := binary.Read(conn, binary.BigEndian, &value); err != nil {
		return 0.0, err
	}

	return value, nil
}

func readInt32(conn *net.UnixConn) (int, error) {
	var format uint8
	if err := binary.Read(conn, binary.BigEndian, &format); err != nil {
		if err == io.EOF {
			return 0, errEndOfInput
		}

		return 0, err
	}

	if format == formatEOF {
		return 0, errEndOfInput
	}

	if format != formatInt32 {
		return 0, errInvalidFormat
	}

	var value uint32
	if err := binary.Read(conn, binary.BigEndian, &value); err != nil {
		return 0, err
	}

	return int(value), nil
}

func sendCommand(conn *net.UnixConn, command string) error {
	if _, err := conn.Write([]byte(command)); err != nil {
		return err
	}

	return nil
}

func closeConnection(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
