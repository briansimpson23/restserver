package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	cfg = make(map[string]string)
)

//func OpenDB() *sql.DB {

func readConfigFile(filename string) map[string]string {

	fd, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not open %q for reading: %s\n", filename, err)
		os.Exit(1)
	}
	defer fd.Close()

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not read %q: %s\n", filename, err)
		os.Exit(1)
	}

	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		if len(line) == 0 || line[0] == 59 {
			continue
		}
		tmp := strings.Split(line, "=")
		cfg[strings.TrimSpace(tmp[0])] = strings.TrimSpace(tmp[1])
	}

	//TODO - need to figure out why this logging doesn't work
	//Debug(fmt.Sprintf("reloaded the config file: %s", filename))
	return cfg
}
