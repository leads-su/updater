package updater

import (
	"encoding/json"
	"fmt"

	"github.com/leads-su/logger"
	"github.com/leads-su/version"
)

// GiteaUpdater defines structure for Gitea Updater
type GiteaUpdater struct {
	UpdaterInterface
	Scheme      string
	Host        string
	Port        uint
	Api         string
	Owner       string
	Repository  string
	BasePath    string
	AccessToken string
}

// GiteaRelease defines structure for each Gitea Release
type GiteaRelease struct {
	Name            string `json:"name"`
	TagName         string `json:"tag_name"`
	CreatedAt       string `json:"created_at"`
	ReleasedAt      string `json:"published_at"`
	UpcomingRelease bool   `json:"prerelease"`
}

// GiteaOptions defines structure for options which can be passed to configurator
type GiteaOptions struct {
	Scheme      string
	Host        string
	Port        uint
	Owner       string
	Repository  string
	AccessToken string
}

// InitializeGitea creates new instance of Gitea updater
func InitializeGitea(options GiteaOptions) (*GiteaUpdater, error) {
	updater := &GiteaUpdater{
		Scheme:      options.Scheme,
		Host:        options.Host,
		Port:        options.Port,
		Api:         "api/v1",
		Owner:       options.Owner,
		Repository:  options.Repository,
		AccessToken: options.AccessToken,
	}

	if updater.Scheme == "" {
		updater.Scheme = "http"
	}

	if updater.Host == "" {
		updater.Host = "gitea.local"
	}

	if updater.Port == 0 {
		updater.Port = 80
	}

	if updater.Owner == "" {
		return nil, fmt.Errorf("you have to specify repository owner name")
	}

	if updater.Repository == "" {
		return nil, fmt.Errorf("you have to specify repository name")
	}

	updater.BasePath = updater.buildBasePath()
	return updater, nil
}

// buildBasePath generates base request path for GitLab
func (updater *GiteaUpdater) buildBasePath() string {
	return fmt.Sprintf("%s://%s:%d/%s/repos/%s/%s",
		updater.Scheme,
		updater.Host,
		updater.Port,
		updater.Api,
		updater.Owner,
		updater.Repository,
	)
}

// Releases will return list of all releases for application
func (updater *GiteaUpdater) Releases() ([]GiteaRelease, error) {
	releases := make([]GiteaRelease, 0)
	requestUrl := fmt.Sprintf("%s/releases", updater.BasePath)
	body, err := sendRequest(requestUrl, "Authorization", "token "+updater.AccessToken)
	if err == nil {
		err = json.Unmarshal(body, &releases)
		if err != nil {
			logger.Warnf("updater:gitlab", "failed to decode releases array - %s", err.Error())
			return nil, err
		}
	}
	return releases, nil
}

// GetLatestVersion returns latest version
func (updater *GiteaUpdater) GetLatestVersion() string {
	currentVersion := version.GetVersion()
	releases, err := updater.Releases()
	latestVersion := currentVersion

	if err == nil {
		for _, release := range releases {
			releaseVersion := release.TagName
			if isRemoteVersionNewer(latestVersion, releaseVersion) {
				latestVersion = releaseVersion
			}
		}
	}
	return latestVersion
}

// IsNewerAvailable checks if currently installed version of CCM is latest
func (updater *GiteaUpdater) IsNewerAvailable() (bool, string, string) {
	currentVersion := version.GetVersion()
	latestVersion := updater.GetLatestVersion()
	return isRemoteVersionNewer(currentVersion, latestVersion), currentVersion, latestVersion
}

// CheckLatest check if latest version of CCM is available
func (updater *GiteaUpdater) CheckLatest() {
	isNewerAvailable, currentVersion, latestVersion := updater.IsNewerAvailable()
	if isNewerAvailable {
		logger.Warnf("updater:gitea", "newer version is available %s, currently running %s", latestVersion, currentVersion)
	}
}

// GiteaCheckLatest check latest version without manually creating new instance
func GiteaCheckLatest(options GiteaOptions) {
	instance, err := InitializeGitea(options)
	if err != nil {
		instance.CheckLatest()
	}
}
