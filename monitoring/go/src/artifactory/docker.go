package artifactory

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/config"
	"github.com/prometheus/common/log"
)

type DockerStats struct {
	DockerRepo string `json:"dockerRepo"`
	Status     int    `json:"status"`
	Slug       string `json:"slug"`
	Error      string `json:"error"`
}

type ArtifactDetails struct {
	RepoUrl string `json:"dockerRepo"`
	Slug    string `json:"slug"`
	Type    string `json:"docker"`
}

type Event struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func VerifyDockerDownloads() ([]DockerStats, error) {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	var artifactDetails = []ArtifactDetails{
		{
			RepoUrl: "infacloud-ct-docker.jfrog.io/ct-busybox1",
			Slug:    "infa_cloud",
			Type:    "docker",
		},
		{
			RepoUrl: "ct-docker.artifacts.cloudtrust.rocks/ct-busybox1",
			Slug:    "ct_docker",
			Type:    "docker",
		},
	}
	dockerStats := make([]DockerStats, len(artifactDetails))

	if err != nil {
		return dockerStats, err
	}

	defer cli.Close()

	conf, conf_err := config.NewConfig()

	if conf_err != nil {
		log.Errorf("Error creating the config. err: %s", conf_err)
		os.Exit(1)
	}

	dockerCredentials := *conf.DockerCredentials

	var authConfig = types.AuthConfig{
		Username:      dockerCredentials.Username,
		Password:      dockerCredentials.Password,
		ServerAddress: "https://infacloud.jfrog.io/v1/",
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	for i := 0; i < len(artifactDetails); i++ {

		dockerStats[i].DockerRepo = artifactDetails[i].RepoUrl
		dockerStats[i].Slug = artifactDetails[i].Slug

		reader, err_d := cli.ImagePull(ctx, artifactDetails[i].RepoUrl, types.ImagePullOptions{RegistryAuth: authConfigEncoded})

		if err_d != nil {
			dockerStats[i].Status = 0
			dockerStats[i].Error = err_d.Error()
			fmt.Println("Error : ", err_d.Error())
		} else {

			defer reader.Close()

			type ErrorMessage struct {
				Error string
			}
			var errorMessage ErrorMessage
			buffIOReader := bufio.NewReader(reader)

			for {
				streamBytes, err := buffIOReader.ReadBytes('\n')
				if err == io.EOF {
					break
				}
				json.Unmarshal(streamBytes, &errorMessage)
				if errorMessage.Error != "" {
					//panic(errorMessage.Error)
					dockerStats[i].Status = 1
					dockerStats[i].Error = errorMessage.Error
					return dockerStats, nil
				}
			}

			io.Copy(os.Stdout, reader)
		}

	}

	return dockerStats, nil

}
