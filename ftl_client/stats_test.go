package ftl_client

import (
	"net"
	"os"
	"reflect"
	"testing"
	"time"
)

const someTimeout = 5 * time.Second

var stats = []byte{0xd2, 0x00, 0x01, 0x72, 0x65, 0xd2, 0x00, 0x00, 0x0c, 0x8e, 0xd2, 0x00, 0x00, 0x00, 0x19, 0xca, 0x3f, 0x47, 0x20, 0xfa, 0xd2, 0x00, 0x00, 0x0e, 0x5f, 0xd2, 0x00, 0x00, 0x01, 0xb3, 0xd2, 0x00, 0x00, 0x0a, 0xc2, 0xd2, 0x00, 0x00, 0x00, 0x07, 0xd2, 0x00, 0x00, 0x00, 0x05, 0xcc, 0x01, 0xc1}

func statsServer(t *testing.T, ln *net.UnixListener) {
	c, err := ln.AcceptUnix()
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 512)
	nr, err := c.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	data := string(buf[0:nr])

	if data != ">stats" {
		t.Errorf("Received unexpected command: %s", data)
	}

	_, err = c.Write(stats)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetStats_missing_address(t *testing.T) {
	client := &FTLClient{
		addr: nil,
	}
	_, err := client.GetStats()
	if err == nil {
		t.Error("GetStats() should fail to connect")
	}
}

func TestGetStats(t *testing.T) {
	socket := testUnixAddr()
	addr, err := net.ResolveUnixAddr("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(socket)

	ln, err := net.ListenUnix("unix", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	ln.Addr()
	ln.SetDeadline(time.Now().Add(someTimeout))

	go statsServer(t, ln)

	client := &FTLClient{
		addr: addr,
	}
	got, err := client.GetStats()
	if err != nil {
		t.Fatal(err)
	}

	want := &Stats{
		DomainsBeingBlocked: 94821,
		DnsQueriesToday:     3214,
		AdsBlockedToday:     25,
		AdsPercentageToday:  0.77784693,
		UniqueDomains:       3679,
		QueriesForwarded:    435,
		QueriesCached:       2754,
		ClientsEverSeen:     7,
		UniqueClients:       5,
		Status:              1,
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("GetDBStats() got = %v, want %v", got, want)
	}
}
