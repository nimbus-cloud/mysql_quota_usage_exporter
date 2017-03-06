package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/nimbus-cloud/mysql_quota_usage_exporter/database"
	"log"
	"fmt"
)

type Exporter struct {
	user string
	pass string
	host string
	brokerDBName string
	port int

	quotaUsage *prometheus.GaugeVec
}

func NewExporter(
	namespace string,
	user string,
	pass string,
	host string,
	brokerDBName string,
	port int,
) *Exporter {
	return &Exporter{
		user: user,
		pass: pass,
		host: host,
		brokerDBName: brokerDBName,
		port: port,
		quotaUsage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "quota_usage",
			Help:      "Used quota ratio",
		},
			[]string{"db_name", "db_size_mb", "db_max_storage_mb"},
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.quotaUsage.Describe(ch)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.quotaUsage.Reset()

	if err := e.collect(); err != nil {
		return
	}

	e.quotaUsage.Collect(ch)
}

func (e *Exporter) collect() error {

	db, err := database.NewConnection(e.user, e.pass, e.host, e.brokerDBName, e.port)
	if db != nil {
		defer db.Close()
	}
	if err != nil {
		log.Fatalf("Error connecting to mysql: %s", err)
	}

	mysqlQuotaUsageRepo := database.NewMysqlQuotaUsageRepo(e.brokerDBName, db)
	dbs, err := mysqlQuotaUsageRepo.All()
	if err != nil {
		log.Fatalf("Error querrying mysql quota usage: %s", err)
	}

	for _, db := range dbs {
		sizeMB := fmt.Sprintf("%.1f", db.SizeMB)
		maxStorageMB := fmt.Sprintf("%.1f", db.MaxStorageMB)
		e.quotaUsage.WithLabelValues(db.Name, sizeMB, maxStorageMB).Set(db.SizeMB/db.MaxStorageMB)
	}

	return nil
}
