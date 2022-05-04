package updater

import (
	"io/ioutil"
	"net/http"

	hashiVersion "github.com/hashicorp/go-version"
	"github.com/leads-su/logger"
)

// UpdaterInterface describes list of methods which should be implemented by service actors
type UpdaterInterface interface {
	CheckLatest()
}

// sendRequest send request to remote endpoint with authorization headers
func sendRequest(url, headerKey, headerValue string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Warnf("updater:create-request", "failed to create new request - %s", err.Error())
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(headerKey, headerValue)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Warnf("updater:send-request", "failed to fetch data from endpoint - %s", err.Error())
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Warnf("updater:process-request", "failed to read data from response body - %s", err.Error())
		return nil, err
	}
	return body, nil
}

// newBaseVersion creates new instance of "empty" version
func newBaseVersion() *hashiVersion.Version {
	v, _ := hashiVersion.NewVersion("0.0.0")
	return v
}

// stringToVersion convert string to version
func stringToVersion(versionString string) *hashiVersion.Version {
	v, err := hashiVersion.NewVersion(versionString)
	if err != nil {
		return newBaseVersion()
	}
	return v
}

// versionToString convert version to string
func versionToString(v *hashiVersion.Version) string {
	return v.String()
}

// isRemoteVersionNewer check if remove version is newer
func isRemoteVersionNewer(local, remote string) bool {
	return stringToVersion(remote).GreaterThan(stringToVersion(local))
}
