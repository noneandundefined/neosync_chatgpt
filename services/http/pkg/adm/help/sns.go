package help

import "encoding/hex"

func HelpConvertToSNS(raw []byte) []string {
	if len(raw)%8 != 0 {
		return nil
	}

	var sns []string
	for i := 0; i < len(raw); i += 8 {
		allZero := true
		for _, b := range raw[i : i+8] {
			if b != 0 {
				allZero = false
				break
			}
		}

		if allZero {
			continue
		}

		sn := hex.EncodeToString(raw[i : i+8])
		sns = append(sns, sn)
	}

	return sns
}
