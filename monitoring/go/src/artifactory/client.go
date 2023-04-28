package artifactory

import (
	"crypto/tls"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/config"
)

// Client represents Artifactory HTTP Client
type Client struct {
	URI               string
	authMethod        string
	cred              config.Credentials
	dockerCredentials config.DockerCredentials
	optionalMetrics   config.OptionalMetrics
	client            *http.Client
	logger            log.Logger
}

// NewClient returns an initialized Artifactory HTTP Client.
func NewClient(conf *config.Config) *Client {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !conf.ArtiSSLVerify}}
	client := &http.Client{
		Timeout:   conf.ArtiTimeout,
		Transport: tr,
	}
	return &Client{
		URI:               conf.ArtiScrapeURI,
		authMethod:        conf.Credentials.AuthMethod,
		dockerCredentials: *conf.DockerCredentials,
		cred:              *conf.Credentials,
		optionalMetrics:   conf.OptionalMetrics,
		client:            client,
		logger:            conf.Logger,
	}
}
