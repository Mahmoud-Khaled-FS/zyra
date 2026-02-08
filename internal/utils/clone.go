package utils

import "maps"

func CloneMap(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}

	dst := make(map[string]string, len(src))
	maps.Copy(dst, src)
	return dst
}
