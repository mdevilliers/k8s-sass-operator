package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/Sirupsen/logrus"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/leaderelection"
	"k8s.io/kubernetes/pkg/client/record"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	masterHost  string
	tlsInsecure bool
	certFile    string
	keyFile     string
	caFile      string
	namespace   string
)

var (
	leaseDuration = 15 * time.Second
	renewDuration = 5 * time.Second
	retryPeriod   = 3 * time.Second
)

func init() {

	flag.StringVar(&masterHost, "master", "", "API Server addr, e.g. 'http://127.0.0.1:8080'. Omit parameter to run in on-cluster mode and utilize the service account token.")
	flag.StringVar(&certFile, "cert-file", "", "Path to public TLS certificate file.")
	flag.StringVar(&keyFile, "key-file", "", "Path to private TLS certificate file.")
	flag.StringVar(&caFile, "ca-file", "", "Path to TLS CA file.")

	namespace = os.Getenv("NAMESPACE")
	if len(namespace) == 0 {
		namespace = "default"
	}
}

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("failed to get hostname: %v", err)
	}

	logrus.Info("sass-operator starting...")
	logrus.Info("hostname", hostname)
	logrus.Info("namespace", namespace)

	leaderelection.RunOrDie(leaderelection.LeaderElectionConfig{
		EndpointsMeta: api.ObjectMeta{
			Namespace: namespace,
			Name:      "sass-operator",
		},
		Client: MustCreateClient(masterHost, tlsInsecure, &restclient.TLSClientConfig{
			CertFile: certFile,
			KeyFile:  keyFile,
			CAFile:   caFile,
		}),
		EventRecorder: &record.FakeRecorder{},
		Identity:      id,
		LeaseDuration: leaseDuration,
		RenewDeadline: renewDuration,
		RetryPeriod:   retryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(stop <-chan struct{}) {
				logrus.Info("ego leader...")
				<-stop
				logrus.Info("no longer the leader.")
			},
			OnStoppedLeading: func() {
				logrus.Fatalf("leader election lost.")
			},
			OnNewLeader: func(identity string) {
				logrus.Info("new leader :", identity)

			},
		},
	})
}

// TODO : pick this apart
// tlsConfig isn't modified inside this function.
// The reason it's a pointer is that it's not necessary to have tlsconfig to create a client.
func MustCreateClient(host string, tlsInsecure bool, tlsConfig *restclient.TLSClientConfig) *unversioned.Client {
	if len(host) == 0 {
		c, err := unversioned.NewInCluster()
		if err != nil {
			panic(err)
		}
		return c
	}
	cfg := &restclient.Config{
		Host:  host,
		QPS:   100,
		Burst: 100,
	}
	hostUrl, err := url.Parse(host)
	if err != nil {
		panic(fmt.Sprintf("error parsing host url %s : %v", host, err))
	}
	if hostUrl.Scheme == "https" {
		cfg.TLSClientConfig = *tlsConfig
		cfg.Insecure = tlsInsecure
	}
	c, err := unversioned.New(cfg)
	if err != nil {
		panic(err)
	}
	return c
}
