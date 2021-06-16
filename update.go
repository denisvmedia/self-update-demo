package main

import (
	"bytes"
	"encoding/json"
	"github.com/hashicorp/go-version"
	"io"
	"net/http"
	"os"
)

type UpdateInfo struct {
	LastVer    string `json:"lastVer"`
	ReleasedAt int    `json:"releasedAt"`
	Sha256     string `json:"sha256"`
}

func fetchUpdateInfo(srv string) (info UpdateInfo, err error) {
	resp, err := http.Get(srv + "/update.json")
	if err != nil {
		return UpdateInfo{}, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return UpdateInfo{}, err
	}

	err = json.Unmarshal(buf.Bytes(), &info)
	if err != nil {
		return UpdateInfo{}, err
	}

	return info, nil
}

func compareVersions(oldVer, newVer string) (bool, error) {
	v1, err := version.NewVersion(oldVer)
	if err != nil {
		return false, err
	}

	v2, err := version.NewVersion(newVer)
	if err != nil {
		return false, err
	}

	return v2.GreaterThan(v1), nil
}

func checkNewVersion(srv, oldVer string) (hasNewVersion bool, info UpdateInfo, err error) {
	info, err = fetchUpdateInfo(srv)
	if err != nil {
		return false, UpdateInfo{}, err
	}

	hasNewVersion, err = compareVersions(oldVer, info.LastVer)
	if err != nil {
		return false, UpdateInfo{}, err
	}

	return hasNewVersion, info, nil
}

func download(f *os.File, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
