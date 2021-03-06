package goappmonitor

import (
	"strings"

	"github.com/wgliang/metrics"
)

func init() {
	// init cfg
	err := loadConfig()
	if err != nil {
		setDefaultConfig()
	}
	cfg := config()

	// init http
	if cfg.Http.Enabled {
		go startHttp(cfg.Http.Listen, cfg.Debug)
	}

	// base collector cron
	if len(cfg.Bases) > 0 {
		go collectBase(cfg.Bases)
	}

	// push cron
	if cfg.Push.Enabled {
		go push2Falcon()
	}
}

// Gauge Actions
func Gauge(name string, value int64) {
	SetGaugeValue(name, float64(value))
}

func GaugeFloat64(name string, value float64) {
	SetGaugeValue(name, value)
}

func SetGaugeValue(name string, value float64) {
	rr := appGaugeFloat64.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.GaugeFloat64); ok {
			r.Update(value)
		}
		return
	}

	r := metrics.NewGaugeFloat64()
	r.Update(value)
	if err := appGaugeFloat64.Collector(name, r); isDuplicateMetricError(err) {
		r := appGaugeFloat64.Get(name).(metrics.GaugeFloat64)
		r.Update(value)
	}
}

func GetGaugeValue(name string) float64 {
	rr := appGaugeFloat64.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.GaugeFloat64); ok {
			return r.Value()
		}
	}
	return 0.0
}

// Meter Actions
func Meter(name string, count int64) {
	SetMeterCount(name, count)
}

func SetMeterCount(name string, count int64) {
	rr := appMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			r.Mark(count)
		}
		return
	}

	r := metrics.NewMeter()
	r.Mark(count)
	if err := appMeter.Collector(name, r); isDuplicateMetricError(err) {
		r := appMeter.Get(name).(metrics.Meter)
		r.Mark(count)
	}
}

func GetMeterCount(name string) int64 {
	rr := appMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.Count()
		}
	}
	return 0
}

func GetMeterRateStep(name string) float64 {
	rr := appMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.RateMean()
		}
	}
	return 0.0
}

func GetMeterRateMean(name string) float64 {
	rr := appMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.RateMean()
		}
	}
	return 0.0
}

func GetMeterRate1(name string) float64 {
	rr := appMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.Rate1()
		}
	}
	return 0.0
}

func GetMeterRate5(name string) float64 {
	rr := appMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.Rate5()
		}
	}
	return 0.0
}

func GetMeterRate15(name string) float64 {
	rr := appMeter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			return r.Rate15()
		}
	}
	return 0.0
}

// Histogram Actions
func Histogram(name string, count int64) {
	SetHistogramCount(name, count)
}

func SetHistogramCount(name string, count int64) {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			r.Update(count)
		}
		return
	}

	s := metrics.NewExpDecaySample(1028, 0.015)
	r := metrics.NewHistogram(s)
	r.Update(count)
	if err := appHistogram.Collector(name, r); isDuplicateMetricError(err) {
		r := appHistogram.Get(name).(metrics.Histogram)
		r.Update(count)
	}
}

func GetHistogramCount(name string) int64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Count()
		}
	}
	return 0
}

func GetHistogramMax(name string) int64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Max()
		}
	}
	return 0
}

func GetHistogramMin(name string) int64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Min()
		}
	}
	return 0
}

func GetHistogramSum(name string) int64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Sum()
		}
	}
	return 0
}

func GetHistogramMean(name string) float64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Mean()
		}
	}
	return 0.0
}

func GetHistogramStdDev(name string) float64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.StdDev()
		}
	}
	return 0.0
}

func GetHistogram50th(name string) float64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.5)
		}
	}
	return 0.0
}

func GetHistogram75th(name string) float64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.75)
		}
	}
	return 0.0
}

func GetHistogram95th(name string) float64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.95)
		}
	}
	return 0.0
}

func GetHistogram99th(name string) float64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.99)
		}
	}
	return 0.0
}

func GetHistogram999th(name string) float64 {
	rr := appHistogram.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Histogram); ok {
			return r.Percentile(0.999)
		}
	}
	return 0.0
}

// Counter Actions
func Counter(name string, count int64) {
	SetCounterCount(name, count)
}

func SetCounterCount(name string, count int64) {
	rr := appCounter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Counter); ok {
			r.Inc(count)
		}
		return
	}

	r := metrics.NewCounter()
	r.Inc(count)
	if err := appCounter.Collector(name, r); isDuplicateMetricError(err) {
		r := appCounter.Get(name).(metrics.Counter)
		r.Inc(count)
	}
}

func GetCounterCount(name string) int64 {
	rr := appCounter.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Counter); ok {
			return r.Count()
		}
	}
	return 0
}

// Self Actions
func selfGauge(name string, value int64) {
	rr := appSelf.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Gauge); ok {
			r.Update(value)
		}
		return
	}

	r := metrics.NewGauge()
	r.Update(value)
	if err := appSelf.Collector(name, r); isDuplicateMetricError(err) {
		r := appSelf.Get(name).(metrics.Gauge)
		r.Update(value)
	}
}

func selfMeter(name string, value int64) {
	rr := appSelf.Get(name)
	if rr != nil {
		if r, ok := rr.(metrics.Meter); ok {
			r.Mark(value)
		}
		return
	}

	r := metrics.NewMeter()
	r.Mark(value)
	if err := appSelf.Collector(name, r); isDuplicateMetricError(err) {
		r := appSelf.Get(name).(metrics.Meter)
		r.Mark(value)
	}
}

// Duplicate Check
func isDuplicateMetricError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Index(err.Error(), "duplicate metric:") == 0
}
