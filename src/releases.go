package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strconv"
	"time"
)

const (
	releaseExp string = `^go(\d+)\.?(\d*)\.?(\d*)$`

	releasesPath string = "https://api.github.com/repos/%s/%s/tags?page=%d"

	releasesOwner string = "golang"
	releasesRepo  string = "go"

	releaseFetchReset int64 = 3600 //Seconds
)

var ()

type goRelease struct {
	Name string

	major int
	minor int
	patch int
}

func parseRelease(str string) *goRelease {
	exp := regexp.MustCompile(releaseExp)

	if !exp.MatchString(str) {
		return nil
	}

	matches := exp.FindStringSubmatch(str)

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		major = -1
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		minor = -1
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		patch = -1
	}

	return &goRelease{
		Name:  matches[0],
		major: major,
		minor: minor,
		patch: patch,
	}
}

type goReleases []goRelease

func parseReleases(releases []string) goReleases {
	res := make(goReleases, 0)

	for _, release := range releases {
		rel := parseRelease(release)

		if rel != nil {
			res = append(res, *rel)
		}
	}

	return res
}

func fetchReleases() (goReleases, error) {
	local := new(releasesJSON)

	err := local.Get()
	if err != nil {
		return nil, err
	}

	if local.Fetched+releaseFetchReset >= time.Now().Unix() {
		return parseReleases(local.Releases), nil
	}

	list := make([]string, 0)

	page := 0
	for {
		resp, err := http.Get(fmt.Sprintf(releasesPath, releasesOwner, releasesRepo, page))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Status code %d received from API", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		relMap := make([]map[string]interface{}, 0)
		err = json.Unmarshal(body, &relMap)
		if err != nil {
			return nil, err
		}

		if len(relMap) == 0 {
			break
		}

		for _, rel := range relMap {
			list = append(list, rel["name"].(string))
		}

		page++
	}

	releases := parseReleases(list)
	releases.sortAsc()

	local.Fetched = time.Now().Unix()
	local.Releases = releases.names()
	local.Save()

	return releases, nil
}

func (releases goReleases) sortAsc() {
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].major < releases[j].major ||
			releases[i].major == releases[j].major && releases[i].minor < releases[j].minor ||
			releases[i].major == releases[j].major && releases[i].minor == releases[j].minor && releases[i].patch < releases[j].patch
	})
}

func (releases goReleases) names() []string {
	names := make([]string, 0, len(releases))

	for _, release := range releases {
		names = append(names, release.Name)
	}

	return names
}

type releasesJSON struct {
	Fetched  int64    `json:"fetched"`
	Releases []string `json:"releases"`
}

func (releases *releasesJSON) Get() error {
	exePath, err := executableDir()
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(path.Join(exePath, "releases.json"))
	if err != nil {
		return err
	}

	return json.Unmarshal(file, releases)
}

func (releases *releasesJSON) Save() error {
	file, err := json.Marshal(releases)
	if err != nil {
		return err
	}

	exePath, err := executableDir()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(exePath, "releases.json"), file, 0644)
}
