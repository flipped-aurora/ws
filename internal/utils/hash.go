package utils

import "hash/crc32"

var table = crc32.MakeTable(1)

func HashUint16(s string) uint16 {
	return uint16(crc32.Checksum([]byte(s), table))
}
