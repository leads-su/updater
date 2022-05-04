package updater

import (
	"encoding/json"
	"fmt"

	"github.com/leads-su/logger"
	"github.com/leads-su/version"
)

// GitlabUpdater defines structure for Gitlab Updater
type GitlabUpdater struct {
	UpdaterInterface
	Scheme      string
	Host        string
	Port        uint
	Api         string
	Projects    string
	ProjectID   uint
	BasePath    string
	AccessToken string
}

// GitlabRelease defines structure for each Gitlab Release
type GitlabRelease struct {
	Name            string `json:"name"`
	TagName         string `json:"tag_name"`
	CreatedAt       string `json:"created_at"`
	ReleasedAt      string `json:"released_at"`
	UpcomingRelease bool   `json:"upcoming_release"`
}

// GitlabOptions defines structure for options which can be passed to configurator
type GitlabOptions struct {
	Scheme      string
	Host        string
	Port        uint
	ApiVersion  uint
	ProjectID   uint
	AccessToken string
}

// InitializeGitlab creates new instance of GitLab updater
func InitializeGitlab(options GitlabOptions) (*GitlabUpdater, error) {
	updater := &GitlabUpdater{
		Scheme:      options.Scheme,
		Host:        options.Host,
		Port:        options.Port,
		Api:         fmt.Sprintf("api/v%d", options.ApiVersion),
		Projects:    "projects",
		ProjectID:   options.ProjectID,
		AccessToken: options.AccessToken,
	}

	if updater.Scheme == "" {
		updater.Scheme = "https"
	}

	if updater.Host == "" {
		updater.Host = "gitlab.com"
	}

	if updater.Port == 0 {
		updater.Port = 443
	}

	if options.ApiVersion == 0 {
		updater.Api = "api/v4"
	}

	if options.ProjectID == 0 {
		return nil, fmt.Errorf("you have to specify ProjectID")
	}

	updater.BasePath = updater.buildBasePath()
	return updater, nil
}

// buildBasePath generates base request path for GitLab
func (updater *GitlabUpdater) buildBasePath() string {
	return fmt.Sprintf("%s://%s:%d/%s/%s/%d",
		updater.Scheme,
		updater.Host,
		updater.Port,
		updater.Api,
		updater.Projects,
		updater.ProjectID,
	)
}

// Releases will return list of all releases for application
func (updater *GitlabUpdater) Releases() ([]GitlabRelease, error) {
	releases := make([]GitlabRelease, 0)
	requestUrl := fmt.Sprintf("%s/releases", updater.BasePath)
	body, err := sendRequest(requestUrl, "PRIVATE-TOKEN", updater.AccessToken)
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
func (updater *GitlabUpdater) GetLatestVersion() string {
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
func (updater *GitlabUpdater) IsNewerAvailable() (bool, string, string) {
	currentVersion := version.GetVersion()
	latestVersion := updater.GetLatestVersion()
	return isRemoteVersionNewer(currentVersion, latestVersion), currentVersion, latestVersion
}

// CheckLatest check if latest version of CCM is available
func (updater *GitlabUpdater) CheckLatest() {
	isNewerAvailable, currentVersion, latestVersion := updater.IsNewerAvailable()
	if isNewerAvailable {
		logger.Warnf("updater:gitlab", "newer version is available %s, currently running %s", latestVersion, currentVersion)
	}
}

// GitlabCheckLatest check latest version without manually creating new instance
func GitlabCheckLatest(options GitlabOptions) {
	instance, err := InitializeGitlab(options)
	if err != nil {
		instance.CheckLatest()
	}
}
