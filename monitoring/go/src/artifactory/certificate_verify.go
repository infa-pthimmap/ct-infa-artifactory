package artifactory

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type CertificateStats struct {
	ArtifcatoryUrl string `json:"artifactoryUrl"`
	DaysToExpire   int    `json:"daysToExpire"`
	ExpiresOn      string `json:"expiresOn"`
	Error          string `json:"error"`
	Status         int    `json:"status"`
}

func GetCertificatesDetails() (CertificateStats, error) {

	host := "infacloud.jfrog.io"
	port := "443"

	certificateStats := CertificateStats{}

	certificateStats.ArtifcatoryUrl = host

	// connect to the website
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		certificateStats.Error = fmt.Sprintf("Error connecting to website: %s", err.Error())
		certificateStats.Status = 0
	}
	defer conn.Close()

	// configure TLS and get the certificate chain
	config := &tls.Config{
		ServerName: host,
	}
	tlsConn := tls.Client(conn, config)
	err = tlsConn.Handshake()
	if err != nil {
		certificateStats.Error = fmt.Sprintf("Error establishing TLS connection: %s", err.Error())
		certificateStats.Status = 0
	}

	defer tlsConn.Close()
	chain := tlsConn.ConnectionState().PeerCertificates

	// get the expiry date of the first certificate in the chain
	expiry := chain[0].NotAfter
	daysLeft := int(expiry.Sub(time.Now()).Hours() / 24)

	// print the expiry date in days from today
	fmt.Println("Certificate expiry date:", expiry.Format("2006-01-02"))
	fmt.Println("Days left:", daysLeft)

	certificateStats.DaysToExpire = daysLeft
	certificateStats.ExpiresOn = fmt.Sprintf("%s", expiry.Format("2006-01-02"))
	certificateStats.Status = 1

	return certificateStats, nil

}
