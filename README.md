# Introduction

This is a small script to parse Polar H10 data and upload it on prometheus to
create a graph with grafana, used with RaspberryPi 3b+

You need also:
- [pushgateway](https://github.com/prometheus/pushgateway)
- [prometheus](https://github.com/prometheus/prometheus)
- [grafana](https://github.com/grafana/grafana)


# Compilation
```sh
go build bpm.go
```

# Usage
```sh
./bpm -b 01:AB:CD:EF:02:AB -j cardiac_frequency -m bpm -u http://192.168.1.2:9091 -v
```

```
Usage of ./bpm:
  -b string
    	Bluetooth mac address used with gatttool to connect and parse data
  -j string
    	Specify prometheus job. (default "cardiac_frequency")
  -m string
    	Specify prometheus metric. (default "bpm")
  -u string
    	pushgateway url to push bpm to prometheus
  -v	verbose
```

# Help
- https://blog.alikhalil.tech/2014/11/polar-h7-bluetooth-le-heart-rate-sensor-on-ubuntu-14-04/
- https://nob.ro/post/polar_h10_ubuntu/
- https://reprage.com/post/how-to-connect-the-raspberry-pi-to-a-bluetooth-heart-rate-monitor
- https://stackoverflow.com/questions/32947807/cannot-connect-to-ble-device-on-raspberry-pi
- https://www.linuxfordevices.com/tutorials/connect-bluetooth-command-line
