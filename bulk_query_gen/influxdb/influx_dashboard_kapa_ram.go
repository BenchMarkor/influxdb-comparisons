package influxdb

import "time"
import (
	"fmt"
	bulkQuerygen "github.com/influxdata/influxdb-comparisons/bulk_query_gen"
)

// InfluxDashboardKapaRam produces Influx-specific queries for the dashboard single-host case.
type InfluxDashboardKapaRam struct {
	InfluxDashboard
	queryTimeRange time.Duration
}

func NewInfluxQLDashboardKapaRam(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(InfluxQL, dbConfig, interval, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardKapaRam{
		InfluxDashboard: *underlying,
		queryTimeRange:  duration,
	}
}

func NewFluxDashboardKapaRam(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(Flux, dbConfig, interval, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardKapaRam{
		InfluxDashboard: *underlying,
		queryTimeRange:  duration,
	}
}

func (d *InfluxDashboardKapaRam) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool

	interval := d.AllInterval.RandWindow(d.queryTimeRange)

	var query string
	//SELECT "used_percent" FROM "telegraf"."autogen"."mem" WHERE time > :dashboardTime: AND "host"='kapacitor'
	query = fmt.Sprintf("SELECT \"used_percent\" FROM system WHERE  hostname='kapacitor' and time >= '%s' and time < '%s'", interval.StartString(), interval.EndString())

	humanLabel := fmt.Sprintf("InfluxDB (%s) kapa mem used in %s", d.language.String(), d.queryTimeRange)

	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
	return q
}
