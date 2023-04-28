package services

import (
	"fmt"

	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/artifactory"
)

type NotifyMetadata struct {
	CertificateStats artifactory.CertificateStats
	HelmStats        []artifactory.HelmStats
	DockerStats      []artifactory.DockerStats
}

type NotificationService struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func SendNotification(notifyMetadata NotifyMetadata) NotificationService {

	//SendEmail("", "")

	fmt.Println(notifyMetadata)

	notificationService := NotificationService{}

	notificationService.Status = 200

	return notificationService

}

func ValidateJfrogServices() NotifyMetadata {

	notifyMetadata := NotifyMetadata{}

	certStats, err_cert := artifactory.GetCertificatesDetails()

	if err_cert != nil {
		fmt.Println(err_cert.Error)
	}

	notifyMetadata.CertificateStats = certStats

	dockerStats, err_docker := artifactory.VerifyDockerDownloads()

	if err_docker != nil {
		fmt.Println(err_docker.Error)
	}

	notifyMetadata.DockerStats = dockerStats

	helmStats, err_helm := artifactory.VerifyHelmArtDownloads()

	if err_helm != nil {
		fmt.Println(err_helm.Error)
	}

	notifyMetadata.HelmStats = helmStats

	return notifyMetadata

}

func ProcessJfrogStatus() NotificationService {

	notifyMetadata := ValidateJfrogServices()

	notificationService := SendNotification(notifyMetadata)

	return notificationService

}

func ProcessJfrogStatusWithReport() NotificationService {

	notifyMetadata := ValidateJfrogServices()

	notificationService := SendNotification(notifyMetadata)

	return notificationService

}
