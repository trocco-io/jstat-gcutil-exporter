# jstat -gcutil Exporter for Prometheus
Exports jstat -gcutil result for Prometheus consumption.

# Help on flags of jstat_exporter
```
  -jstat.path string
    	jstat path (default "/usr/bin/jstat")
  -target.pid string
    	target pid (default ":0")
  -web.listen-address string
    	Address on which to expose metrics and web interface. (default ":9010")
  -web.telemetry-path string
    	Path under which to expose metrics. (default "/metrics")
```
