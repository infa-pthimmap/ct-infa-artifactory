package artifactory

import (
	"fmt"
	"os"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/auth"
	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/services"
	jfrog_config "github.com/jfrog/jfrog-client-go/config"
	art_exp_config "github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/config"
	"github.com/prometheus/common/log"
)

type HelmStats struct {
	HelmRepo string `json:"helmRepo"`
	Status   int    `json:"status"`
	Slug     string `json:"slug"`
	Error    string `json:"error"`
}

type ArtifactHelmDetails struct {
	HelmUrl string `json:"helmRepo"`
	Slug    string `json:"slug"`
	Type    string `json:"helm"`
}

func VerifyHelmArtDownloads() ([]HelmStats, error) {

	var artifactDetails = []ArtifactHelmDetails{
		{
			HelmUrl: "helm-virt/email-wrapper1/helm_ew-1.0.tgz",
			Slug:    "helm_virt",
			Type:    "helm",
		},
	}
	helmStats := make([]HelmStats, len(artifactDetails))

	conf, conf_err := art_exp_config.NewConfig()

	if conf_err != nil {
		log.Errorf("Error creating the config. err: %s", conf_err)
		os.Exit(1)
	}

	dockerCredentials := *conf.DockerCredentials

	rtDetails := auth.NewArtifactoryDetails()
	rtDetails.SetUrl("https://infacloud.jfrog.io/artifactory/")
	rtDetails.SetUser(dockerCredentials.Username)
	rtDetails.SetPassword(dockerCredentials.Password)

	serviceConfig, err := jfrog_config.NewConfigBuilder().
		SetServiceDetails(rtDetails).
		Build()

	rtManager, err := artifactory.New(serviceConfig)

	for i := 0; i < len(artifactDetails); i++ {

		helmStats[i].HelmRepo = artifactDetails[i].HelmUrl
		helmStats[i].Slug = artifactDetails[i].Slug

		if err != nil {

			helmStats[i].Status = 0
			helmStats[i].Error = err.Error()

		} else {

			params := services.NewDownloadParams()
			params.Pattern = artifactDetails[i].HelmUrl
			params.Target = "/app/"

			totalDownloaded, totalFailed, err := rtManager.DownloadFiles(params)

			if err != nil {
				helmStats[i].Status = 0
				helmStats[i].Error = err.Error()
			} else {

				if totalDownloaded > 0 {
					helmStats[i].Status = 1
					helmStats[i].Error = ""
				} else {
					helmStats[i].Status = 0
					msg := fmt.Sprintf("%s %s %d", "Failed to download artifact ", ", found downlaod as 0 and totalFailed as ", totalFailed)
					helmStats[i].Error = msg
				}

			}

		}

	}

	return helmStats, nil

}
