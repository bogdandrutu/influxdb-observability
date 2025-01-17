package influx2otel_test

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"testing"
	"time"

	"github.com/influxdata/influxdb-observability/common"
	"github.com/influxdata/influxdb-observability/influx2otel"
	"github.com/stretchr/testify/require"
)

func TestUnknownSchema(t *testing.T) {
	c, err := influx2otel.NewLineProtocolToOtelMetrics(new(common.NoopLogger))
	require.NoError(t, err)

	b := c.NewBatch()
	err = b.AddPoint("cpu",
		map[string]string{
			"container.name":       "42",
			"otel.library.name":    "My Library",
			"otel.library.version": "latest",
			"cpu":                  "cpu4",
			"host":                 "777348dc6343",
		},
		map[string]interface{}{
			"usage_user":   0.10090817356207936,
			"usage_system": 0.3027245206862381,
			"some_int_key": int64(7),
		},
		time.Unix(0, 1395066363000000123),
		common.InfluxMetricValueTypeUntyped)
	require.NoError(t, err)

	expect := pmetric.NewMetrics()
	rm := expect.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().InsertString("container.name", "42")
	ilMetrics := rm.ScopeMetrics().AppendEmpty()
	ilMetrics.Scope().SetName("My Library")
	ilMetrics.Scope().SetVersion("latest")
	m := ilMetrics.Metrics().AppendEmpty()
	m.SetName("cpu_usage_user")
	m.SetDataType(pmetric.MetricDataTypeGauge)
	dp := m.Gauge().DataPoints().AppendEmpty()
	dp.Attributes().InsertString("cpu", "cpu4")
	dp.Attributes().InsertString("host", "777348dc6343")
	dp.SetTimestamp(pcommon.Timestamp(1395066363000000123))
	dp.SetDoubleVal(0.10090817356207936)
	m = ilMetrics.Metrics().AppendEmpty()
	m.SetName("cpu_usage_system")
	m.SetDataType(pmetric.MetricDataTypeGauge)
	dp = m.Gauge().DataPoints().AppendEmpty()
	dp.Attributes().InsertString("cpu", "cpu4")
	dp.Attributes().InsertString("host", "777348dc6343")
	dp.SetTimestamp(pcommon.Timestamp(1395066363000000123))
	dp.SetDoubleVal(0.3027245206862381)
	m = ilMetrics.Metrics().AppendEmpty()
	m.SetName("cpu_some_int_key")
	m.SetDataType(pmetric.MetricDataTypeGauge)
	dp = m.Gauge().DataPoints().AppendEmpty()
	dp.Attributes().InsertString("cpu", "cpu4")
	dp.Attributes().InsertString("host", "777348dc6343")
	dp.SetTimestamp(pcommon.Timestamp(1395066363000000123))
	dp.SetIntVal(7)

	assertMetricsEqual(t, expect, b.GetMetrics())
}
