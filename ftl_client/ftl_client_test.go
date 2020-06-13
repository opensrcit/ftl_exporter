package ftl_client

import (
	"io/ioutil"
	"os"
)

// testUnixAddr uses ioutil.TempFile to get a name that is unique.
func testUnixAddr() string {
	f, err := ioutil.TempFile("", "ftl_client_test")
	if err != nil {
		panic(err)
	}
	addr := f.Name()
	f.Close()
	os.Remove(addr)

	return addr
}
