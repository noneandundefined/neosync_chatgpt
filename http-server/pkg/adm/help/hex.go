package help

import "encoding/hex"

func HelpConvertToHex(raw []byte) string {
	return hex.EncodeToString(raw)
}
