package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var lc int
var freq int
var m map[string]*Stats

type Stats struct {
	total, errors int
}

func (stats *Stats) Percentage() float64 {
	e := float64(stats.errors)
	t := float64(stats.total)
	p := (e / t) * 100.0
	return p
}

func Reporter(m map[string]*Stats) {
	for k, v := range m {
		fmt.Printf("%s returned %.2f%% 5xx errors\n", k, v.Percentage())
	}
}

func LogParser(log string) {
	file, err := os.Open(log)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		LineParser(scanner.Text())
	}

}
func HostChecker(hostname string) {
	if _, ok := m[hostname]; !ok {
		m[hostname] = &Stats{
			0, 0,
		}
	}
}

func HttpCodeChecker(httpcode string) int {
	if strings.HasPrefix(httpcode, "5") {
		return 1
	}
	return 0
}

func LineParser(line string) {
	lc += 1
	if lc%freq == 0 {
		fmt.Printf("lines parsed: %d\r", lc)
	}
	values := strings.Split(line, "|")
	hostname := values[2]
	HostChecker(hostname)
	m[hostname].total += 1
	httpcode := values[4]
	m[hostname].errors += HttpCodeChecker(httpcode)
}

func ArgParser() []string {
	flag.Parse()
	files := flag.Args()
	return files
}

func main() {
	lc = 0
	freq = 79351
	logFiles := ArgParser()
	m = make(map[string]*Stats)
	for _, log := range logFiles {
		LogParser(log)
	}
	fmt.Printf("total lines: %d\n", lc)
	Reporter(m)
}
