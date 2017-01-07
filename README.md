# sensorlog

This is a very simple tool to log the output of `senors` (a Linux tool) to a CSV file.

# Usage

First, you must have [Go](https://golang.org/doc/install). Next, download like so:

```
$ go get github.com/unixpickle/sensorlog
```

Now run the command as follows:

```
$ sensorlog -out sensors.csv
```

You can use the `-stripunits` flag to remove units from CSV entries. You can use the `-interval` flag to control how frequently the command runs. See `-help` for more info.
