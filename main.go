package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/molu8bits/s3bucket-exporter/controllers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/Sirupsen/logrus"
)

var (
	up                          = prometheus.NewDesc("s3_endpoint_up", "Conection to S3 successful", []string{"s3name"}, nil)
	listenPort                  = ":9655"
	s3Endpoint                  = ""
	s3AccessKey                 = ""
	s3SecretKey                 = ""
	s3DisableSSL                = false
	s3Name                      = ""
	s3DisableEndpointHostPrefix = false
	s3ForcePathStyle            = true
	s3Region                    = "default"
	s3Conn                      controllers.S3Conn
	logLevel                    = "info"
)

func envString(key, def string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return def
}

func envBool(key string, def bool) bool {
	def2, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return def
	}
	return def2
}

func init() {
	var s3Conn controllers.S3Conn
	flag.StringVar(&s3Endpoint, "s3_endpoint", envString("S3_ENDPOINT", s3Conn.S3ConnEndpoint), "S3_ENDPOINT - eg. myceph.com:7480")
	flag.StringVar(&s3AccessKey, "s3_access_key", envString("S3_ACCESS_KEY", s3AccessKey), "S3_ACCESS_KEY - aws_access_key")
	flag.StringVar(&s3SecretKey, "s3_secret_key", envString("S3_SECRET_KEY", s3SecretKey), "S3_SECRET_KEY - aws_secret_key")
	flag.StringVar(&s3Name, "s3_name", envString("S3_NAME", s3Name), "S3_NAME")
	flag.StringVar(&s3Region, "s3_region", envString("S3_REGION", s3Region), "S3_REGION")
	flag.StringVar(&listenPort, "listen_port", envString("LISTEN_PORT", listenPort), "LISTEN_PORT e.g ':9655'")
	flag.StringVar(&logLevel, "log_level", envString("LOG_LEVEL", logLevel), "LOG_LEVEL")
	flag.BoolVar(&s3DisableSSL, "s3_disable_ssl", envBool("S3_DISABLE_SSL", s3DisableSSL), "s3 disable ssl")
	flag.BoolVar(&s3DisableEndpointHostPrefix, "s3_disable_endpoint_host_prefix", envBool("S3_DISABLE_ENDPOINT_HOST_PREFIX", s3DisableEndpointHostPrefix), "S3_DISABLE_ENDPOINT_HOST_PREFIX")
	flag.BoolVar(&s3ForcePathStyle, "s3_force_path_style", envBool("S3_FORCE_PATH_STYLE", s3ForcePathStyle), "S3_FORCE_PATH_STYLE")
	flag.Parse()
}

// S3Collector dummy struct
type S3Collector struct {
}

// Describe - Implements prometheus.Collector
func (c S3Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

// Collect - Implements prometheus.Collector.
func (c S3Collector) Collect(ch chan<- prometheus.Metric) {
	s3Conn.S3ConnEndpoint = s3Endpoint
	s3Conn.S3ConnAccessKey = s3AccessKey
	s3Conn.S3ConnSecretKey = s3SecretKey
	s3Conn.S3ConnDisableSsl = s3DisableSSL
	s3Conn.S3ConnName = s3Name
	s3Conn.S3ConnDisableEdnpointHostPrefix = s3DisableEndpointHostPrefix
	s3Conn.S3ConnForcePathStyle = s3ForcePathStyle
	s3Conn.S3ConnRegion = s3Region

	s3metrics, err := controllers.S3UsageInfo(s3Conn)
	s3name := s3metrics.S3Name
	if err != nil {
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		return
	}
	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1, s3Name)
	log.Debug("s3metrics read from s3_endpoint :", s3metrics)

	totalObjectNumber := s3metrics.S3ObjectNumber
	totalSize := s3metrics.S3Size
	descS := prometheus.NewDesc("s3_total_size", "S3 Total Bucket Size", []string{"s3name"}, nil)
	descON := prometheus.NewDesc("s3_total_object_number", "S3 Total Object Number", []string{"s3name"}, nil)
	ch <- prometheus.MustNewConstMetric(
		descS, prometheus.GaugeValue, float64(totalSize), s3name)
	ch <- prometheus.MustNewConstMetric(
		descON, prometheus.GaugeValue, float64(totalObjectNumber), s3name)

	for _, s := range s3metrics.S3Buckets {
		bucketName := s.BucketName
		bucketObjectNumber := s.BucketObjectNumber
		bucketSize := s.BucketSize
		descBucketS := prometheus.NewDesc("s3_bucket_size", "S3 metric Total Bucket Size", []string{"s3name", "bucketname"}, nil)
		descBucketON := prometheus.NewDesc("s3_bucket_object_number", "S3 metric Total Object Number", []string{"s3name", "bucketname"}, nil)
		ch <- prometheus.MustNewConstMetric(
			descBucketS, prometheus.GaugeValue, float64(bucketSize), s3name, bucketName)
		ch <- prometheus.MustNewConstMetric(
			descBucketON, prometheus.GaugeValue, float64(bucketObjectNumber), s3name, bucketName)
	}
}

func main() {
	logrusLogLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Println("Unable to set logging level:", logLevel)
		os.Exit(1)
	}
	log.SetLevel(logrusLogLevel)

	if s3Endpoint == "" {
		log.Fatal("s3 endpoint must be configured")
	} else if s3AccessKey == "" {
		log.Fatal("S3_ACCESS_KEY must be configured")
	} else if s3SecretKey == "" {
		log.Fatal("S3_SECRET_KEY must be configured")
	} else if s3Name == "" {
		log.Fatal("S3_NAME must be configured")
	}

	c := S3Collector{}
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port", listenPort)
	log.Info("s3 name '", s3Name, "' available at s3 endpoint '", s3Endpoint, "' will be monitored")
	log.Info("listenPort :", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, nil))

}
