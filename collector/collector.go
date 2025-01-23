package collector

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

const (
	namespace = "ses"
)

// Exporter collects metrics from aws ses.
type Exporter struct {
	max24hoursend    *prometheus.Desc
	maxsendrate      *prometheus.Desc
	sentlast24hours  *prometheus.Desc
	Bounces          *prometheus.Desc
	Complaints       *prometheus.Desc
	DeliveryAttempts *prometheus.Desc
	Rejects          *prometheus.Desc
}

// NewExporter returns an initialized exporter.
func NewExporter() *Exporter {
	return &Exporter{
		max24hoursend: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max24hoursend"),
			"The maximum number of emails allowed to be sent in a rolling 24 hours.",
			[]string{"aws_region"},
			nil,
		),
		maxsendrate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "maxsendrate"),
			"The maximum rate of emails allowed to be sent per second.",
			[]string{"aws_region"},
			nil,
		),
		sentlast24hours: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "sentlast24hours"),
			"The number of emails sent in the last 24 hours.",
			[]string{"aws_region"},
			nil,
		),
		Bounces: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "Bounces"),
			"The number of emails of emails that have bounced.",
			[]string{"aws_region"},
			nil,
		),
		Complaints: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "Complaints"),
			"Number of unwanted emails that were rejected by recipients.",
			[]string{"aws_region"},
			nil,
		),
		DeliveryAttempts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "DeliveryAttempts"),
			"Number of emails that have been sent.",
			[]string{"aws_region"},
			nil,
		),
		Rejects: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "Rejects"),
			"Number of emails rejected by Amazon SES.",
			[]string{"aws_region"},
			nil,
		),
	}
}

// Sort SendDataPoints base on Timestamp to
// retrieve data in last 15 min.
type sortedDataPoint struct {
	Data []*ses.SendDataPoint
}

func (sd sortedDataPoint) Len() int {
	return len(sd.Data)
}

func (sd sortedDataPoint) Less(i, j int) bool {
	return sd.Data[i].Timestamp.Before(*sd.Data[j].Timestamp)
}

func (sd sortedDataPoint) Swap(i, j int) {
	sd.Data[i], sd.Data[j] = sd.Data[j], sd.Data[i]
}

// Describe all the metrics.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.max24hoursend
	ch <- e.maxsendrate
	ch <- e.sentlast24hours
	ch <- e.Bounces
	ch <- e.Complaints
	ch <- e.DeliveryAttempts
	ch <- e.Rejects
}

// Collect fetches the statistics from aws ses sdk, and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	svc := ses.New(session.New())
	regionName := svc.SigningRegion
	sendStatisticsInput := &ses.GetSendStatisticsInput{}
	sendStatisticsOutput, err := svc.GetSendStatistics(sendStatisticsInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Fatalf("Failed to get sending quota from region %s: %s", regionName, aerr)
			}
		} else {
			log.Fatalf("Failed to get sending quota from region %s: %s", regionName, err)
		}
		return
	}

	sortedSendStatistics := sortedDataPoint{
		Data: sendStatisticsOutput.SendDataPoints,
	}
	sort.Sort(sortedSendStatistics)
	latestSendStatistics := sortedSendStatistics.Data[len(sortedSendStatistics.Data)-1]

	sendQuotaInput := &ses.GetSendQuotaInput{}
	sendQuotaOutput, err := svc.GetSendQuota(sendQuotaInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Fatalf("Failed to get sending quota from region %s: %s", regionName, aerr)
			}
		} else {
			log.Fatalf("Failed to get sending quota from region %s: %s", regionName, err)
		}
		return
	}

	ch <- prometheus.MustNewConstMetric(e.max24hoursend, prometheus.GaugeValue, *sendQuotaOutput.Max24HourSend, regionName)
	ch <- prometheus.MustNewConstMetric(e.maxsendrate, prometheus.GaugeValue, *sendQuotaOutput.MaxSendRate, regionName)
	ch <- prometheus.MustNewConstMetric(e.sentlast24hours, prometheus.GaugeValue, *sendQuotaOutput.SentLast24Hours, regionName)
	ch <- prometheus.MustNewConstMetric(e.Bounces, prometheus.GaugeValue, float64(*latestSendStatistics.Bounces), regionName)
	ch <- prometheus.MustNewConstMetric(e.Complaints, prometheus.GaugeValue, float64(*latestSendStatistics.Complaints), regionName)
	ch <- prometheus.MustNewConstMetric(e.DeliveryAttempts, prometheus.GaugeValue, float64(*latestSendStatistics.DeliveryAttempts), regionName)
	ch <- prometheus.MustNewConstMetric(e.Rejects, prometheus.GaugeValue, float64(*latestSendStatistics.Rejects), regionName)
}
