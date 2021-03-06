package influxdb

import "time"
import (
	"fmt"
	bulkQuerygen "github.com/influxdata/influxdb-comparisons/bulk_query_gen"
)

// InfluxDashboardKapaCpu produces Influx-specific queries for the dashboard single-host case.
type InfluxDashboardKapaCpu struct {
	InfluxDashboard
	queryTimeRange time.Duration
}

func NewInfluxQLDashboardKapaCpu(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(InfluxQL, dbConfig, interval, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardKapaCpu{
		InfluxDashboard: *underlying,
		queryTimeRange:  duration,
	}
}

func NewFluxDashboardKapaCpu(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {
	underlying := newInfluxDashboard(Flux, dbConfig, interval, scaleVar).(*InfluxDashboard)
	return &InfluxDashboardKapaCpu{
		InfluxDashboard: *underlying,
		queryTimeRange:  duration,
	}
}

func (d *InfluxDashboardKapaCpu) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool

	interval := d.AllInterval.RandWindow(d.queryTimeRange)

	var query string
	//SELECT 100 - "usage_idle" FROM "telegraf"."autogen"."cpu" WHERE time > now() - 15m AND "cpu"='cpu-total' AND "host"='kapacitor'
	query = fmt.Sprintf("SELECT 100 - \"usage_idle\" FROM cpu WHERE hostname='kapacitor' and time >= '%s' and time < '%s'", interval.StartString(), interval.EndString())

	humanLabel := fmt.Sprintf("InfluxDB (%s) kapa cpu in %s", d.language.String(), d.queryTimeRange)

	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
	return q
}
