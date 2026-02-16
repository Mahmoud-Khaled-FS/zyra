package utils

import "net/url"

func IsValidURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func JoinURL(baseStr, pathStr string) (string, error) {
	base, err := url.Parse(baseStr)
	if err != nil {
		return "", err
	}

	ref, err := url.Parse(pathStr)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(ref).String(), nil
}
