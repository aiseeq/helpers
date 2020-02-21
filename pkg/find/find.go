package find

import (
	"crypto/md5"
	"encoding/hex"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringInMap(a string, list map[string]string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func KeyByValue(arr []string, val string) int {
	for k, v := range arr {
		if v == val {
			return k
		}
	}
	return -1
}

func MapKeyByValue(amap map[string]string, val string) string {
	for key, v := range amap {
		if val == v {
			return key
		}
	}
	return ""
}

func SliceDiff(a, b []string) (ab []string) {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return
}

func MinFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func MaxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func MD5Hash(text []byte) string {
	hasher := md5.New()
	hasher.Write(text)
	return hex.EncodeToString(hasher.Sum(nil))
}
