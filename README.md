# Updater Package for Go Lang
This package provides and easy way to check for application updates.
This package depends on another package - [version](http://github.com/leads-su/version)

# Configuring Update Source
As of now, Gitlab and Gitea are the only implemented update sources supported by the updater

## GitLab
To initialize updater, you can call it in the following manner:
```go
updater.InitializeGitlab(updater.GitlabOptions{
    Scheme:      "https",                // OPTIONAL parameter. Defaults to - https
    Host:        "gitlab.com",           // OPTIONAL parameter. Defaults to - gitlab.com
    Port:        443,                    // OPTIONAL parameter. Defaults to - 443
    ApiVersion:  4,                      // OPTIONAL parameter. Defaults to - 4
    ProjectID:   123,                    // REQUIRED parameter
	  AccessToken: "",                     // REQUIRED parameter
})
```

## Gitea
To initialize updater, you can call it in the following manner:
```go
updater.InitializeGitea(updater.GiteaOptions{
    Scheme      "http",                // OPTIONAL parameter. Defaults to - http
    Host        "gitea.local",         // OPTIONAL parameter. Defaults to - gitea.local
    Port        80,                    // OPTIONAL parameter. Defaults to - 80
    Owner       "owner",               // REQUIRED parameter
    Repository  "repository",          // REQUIRED parameter
    AccessToken "",                    // REQUIRED parameter
})
```

# Providing LDFLAGS
This is the list of LDFLAGS you can provide to the builder while building target binary
 - `-X github.com/leads-su/version.version` - sets version of binary
 - `-X github.com/leads-su/version.commit` - sets commit with which binary was built
 - `-X github.com/leads-su/version.buildDate` - sets build date for binary
 - `-X github.com/leads-su/version.builtBy` - sets builder name for binary

# LDFLAGS in GitLab pipeline
 - `-X github.com/leads-su/version.version={{.Version}}`
 - `-X github.com/leads-su/version.commit={{.Commit}}`
 - `-X github.com/leads-su/version.buildDate={{.Date}}`
 - `-X github.com/leads-su/version.builtBy=builder-name`

# Using GoReleaser (Preferred way of building binaries)
This is an example configuration for GoReleaser (`.goreleaser.yaml`)
## Base for GoReleaser Configuration
```yaml
builds:
  - env:
      - CGO_ENABLED=0
      - GO111MODULE=on
    goos:
      - linux
      - windows
      - darwin
    goarch:
      - 386
      - amd64
    ldflags:
      - -s -w -X github.com/leads-su/version.version={{.Version}} -X github.com/leads-su/version.commit={{.Commit}} -X github.com/leads-su/version.buildDate={{.Date}} -X github.com/leads-su/version.builtBy=goreleaser
archives:
  - format: zip
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    files:
      - changelog*
      - CHANGELOG*
      - readme*
      - README*
checksum:
  name_template: "{{ .ProjectName }}_{{ .Version }}_SHA256SUMS"
  algorithm: sha256
changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
```

## Gitlab Specific Options
These are the options if you want to use Gitlab as your release manager
```yaml
gitlab_urls:
  api: https://gitlab.com/api/v4/
  download: https://gitlab.com
  skip_tls_verify: true
```

## Gitea Specific Options
These are the options if you want to use Gitea as your release manager
```yaml
gitea_urls:
  api: http://gitea.locall/api/v1/
  download: http://gitea.local
  skip_tls_verify: true
```

