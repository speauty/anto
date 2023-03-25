package util

func HasUTF8Dom(bytes []byte) bool {
	return len(bytes) >= 3 && bytes[0] == 0xEF && bytes[1] == 0xBB && bytes[2] == 0xBF
}
