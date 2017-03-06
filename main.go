package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/go-sql-driver/mysql"
)

var (
	listenAddress = flag.String("web.listen", ":9118", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	namespace     = flag.String("namespace", "mysql", "Namespace for the Mysql quota usage metrics.")

	dbUser        = flag.String("db.User", "root", "User to connect to Mysql broker data base.")
	dbPass        = flag.String("db.Pass", "", "User to connect to Mysql broker data base.")
	dbHost        = flag.String("db.Host", "", "User to connect to Mysql broker data base.")
	dbName        = flag.String("db.Name", "mysql_broker", "User to connect to Mysql broker data base.")
	dbPort        = flag.Int("db.Port", 3306, "User to connect to Mysql broker data base.")
)

func main() {
	flag.Parse()

	prometheus.MustRegister(NewExporter(
		*namespace,
		*dbUser,
		*dbPass,
		*dbHost,
		*dbName,
		*dbPort))

	log.Printf("Starting Server: %s", *listenAddress)
	handler := prometheus.Handler()
	if *metricsPath == "" || *metricsPath == "/" {
		http.Handle(*metricsPath, handler)
	} else {
		http.Handle(*metricsPath, handler)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>MySql Quota Usage Exporter</title></head>
			<body>
			<h1>MySql Quota Usage Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		})
	}

	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
