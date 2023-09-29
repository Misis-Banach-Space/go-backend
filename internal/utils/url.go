package utils

import (
	"net/url"
	"strings"
)

func GetUrlDomain(rawUrl string) (string, error) {
	url, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	domain := strings.TrimPrefix(url.Hostname(), "www.")

	prefix := strings.Split(rawUrl, "://")[0]

	return prefix + "://" + domain, nil
}
