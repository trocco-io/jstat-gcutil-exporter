package main

import (
	"flag"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/log"
)

const (
	namespace = "jstat"
)

var (
	listenAddress = flag.String("listen-address", ":9010", "Address on which to expose metrics")
	metricsPath   = flag.String("metrics-path", "/metrics", "Path under which to expose metrics.")
	jstatPath     = flag.String("jstat-path", "/usr/bin/jstat", "jstat path")
	targetPid     = flag.String("pid", ":0", "target pid")
)

type Exporter struct {
	jstatPath  string
	targetPid  string
	s0         prometheus.Gauge
	s1         prometheus.Gauge
	eden       prometheus.Gauge
	old        prometheus.Gauge
	meta       prometheus.Gauge
	ccs        prometheus.Gauge
	ygc        prometheus.Gauge
	ygct       prometheus.Gauge
	fgc        prometheus.Gauge
	fgct       prometheus.Gauge
	gct        prometheus.Gauge
}

func NewExporter(jstatPath string, targetPid string) *Exporter {
	return &Exporter{
		jstatPath: jstatPath,
		targetPid: targetPid,
		s0: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "s0",
			Help:      "s0",
		}),
		s1: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "s1",
			Help:      "s1",
		}),
		eden: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "eden",
			Help:      "eden",
		}),
		old: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "old",
			Help:      "old",
		}),
		meta: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "meta",
			Help:      "meta",
		}),
		ccs: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ccs",
			Help:      "ccs",
		}),
		ygc: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ygc",
			Help:      "ygc",
		}),
		ygct: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "ygct",
			Help:      "ygct",
		}),
		fgc: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "fgc",
			Help:      "fgc",
		}),
		fgct: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "fgct",
			Help:      "fgct",
		}),
		gct: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "gct",
			Help:      "gct",
		}),
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.s0.Describe(ch)
	e.s1.Describe(ch)
	e.eden.Describe(ch)
	e.old.Describe(ch)
	e.meta.Describe(ch)
	e.ccs.Describe(ch)
	e.ygc.Describe(ch)
	e.ygct.Describe(ch)
	e.fgc.Describe(ch)
	e.fgct.Describe(ch)
	e.gct.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.JstatUtil(ch)
}

func (e *Exporter) JstatUtil(ch chan<- prometheus.Metric) {

	out, err := exec.Command(e.jstatPath, "-gcutil", e.targetPid).Output()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range strings.Split(string(out), "\n") {
		if i == 1 {
			parts := strings.Fields(line)

			s0, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.s0.Set(s0)
			e.s0.Collect(ch)

			s1, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.s1.Set(s1)
			e.s1.Collect(ch)

			eden, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.eden.Set(eden)
			e.eden.Collect(ch)

			old, err := strconv.ParseFloat(parts[3], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.old.Set(old)
			e.old.Collect(ch)

			meta, err := strconv.ParseFloat(parts[4], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.meta.Set(meta)
			e.meta.Collect(ch)

			ccs, err := strconv.ParseFloat(parts[4], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.ccs.Set(ccs)
			e.ccs.Collect(ch)

			ygc, err := strconv.ParseFloat(parts[5], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.ygc.Set(ygc)
			e.ygc.Collect(ch)

			ygct, err := strconv.ParseFloat(parts[6], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.ygct.Set(ygct)
			e.ygct.Collect(ch)

			fgc, err := strconv.ParseFloat(parts[7], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.fgc.Set(fgc)
			e.fgc.Collect(ch)

			fgct, err := strconv.ParseFloat(parts[8], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.fgct.Set(fgct)
			e.fgct.Collect(ch)

			gct, err := strconv.ParseFloat(parts[9], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.gct.Set(gct)
			e.gct.Collect(ch)
		}
	}
}

func main() {
	flag.Parse()

	exporter := NewExporter(*jstatPath, *targetPid)
	prometheus.MustRegister(exporter)

	log.Printf("Starting Server: %s", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>jstat Exporter</title></head>
		<body>
		<h1>jstat Exporter</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		</body>
		</html>`))
	})
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

}
