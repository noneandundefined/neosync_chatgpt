package util

import "hash/fnv"

/* Алгоритм хеширования по fnv1a */
func FNV1aHash(data []byte) uint32 {
	h := fnv.New32a()
	h.Write(data)
	return h.Sum32()
}
