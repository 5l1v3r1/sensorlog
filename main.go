package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

func main() {
	var interval time.Duration
	var outFile string
	var stripUnits bool
	flag.DurationVar(&interval, "interval", time.Minute, "sample duration")
	flag.StringVar(&outFile, "out", "sensors.csv", "output file")
	flag.BoolVar(&stripUnits, "stripunits", false, "strip units")
	flag.Parse()

	file, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "open output:", err)
		os.Exit(1)
	}
	defer file.Close()
	off, err := file.Seek(0, os.SEEK_END)
	if err != nil {
		fmt.Fprintln(os.Stderr, "seek:", err)
		os.Exit(1)
	}
	first := off == 0
	writer := csv.NewWriter(file)

	ticker := time.NewTicker(interval)
	for {
		data, err := ReadSensors()
		if err != nil {
			fmt.Fprintln(os.Stderr, "read sensors:", err)
			os.Exit(1)
		}
		writeLine(writer, data, first, stripUnits)
		first = false
		<-ticker.C
	}
}

func writeLine(w *csv.Writer, data map[string]string, first, stripUnits bool) {
	data["time"] = fmt.Sprintf("%d", time.Now().Unix())
	sorter := NewFieldSorter(data)
	sort.Sort(sorter)
	if first {
		w.Write(sorter.Fields)
	}
	if stripUnits {
		for i, x := range sorter.Values {
			sorter.Values[i] = StripUnits(x)
		}
	}
	w.Write(sorter.Values)
	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Fprintln(os.Stderr, "write log entry:", err)
		os.Exit(1)
	}
}

type FieldSorter struct {
	Fields []string
	Values []string
}

func NewFieldSorter(data map[string]string) *FieldSorter {
	res := &FieldSorter{}
	for key, value := range data {
		res.Fields = append(res.Fields, key)
		res.Values = append(res.Values, value)
	}
	return res
}

func (f *FieldSorter) Len() int {
	return len(f.Fields)
}

func (f *FieldSorter) Less(i, j int) bool {
	return f.Fields[i] < f.Fields[j]
}

func (f *FieldSorter) Swap(i, j int) {
	f.Fields[i], f.Fields[j] = f.Fields[j], f.Fields[i]
	f.Values[i], f.Values[j] = f.Values[j], f.Values[i]
}
