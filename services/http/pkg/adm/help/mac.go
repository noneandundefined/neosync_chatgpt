package help

import "fmt"

func HelpConvertToMacArray(raw []byte) []string {
	var macs []string

	for i := 0; i+6 < len(raw); i += 6 {
		chunk := raw[i : i+6]

		empty := true
		for _, b := range chunk {
			if b != 0 {
				empty = false
				break
			}
		}

		if empty {
			continue
		}

		mac := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", chunk[0], chunk[1], chunk[2], chunk[3], chunk[4], chunk[5])
		macs = append(macs, mac)
	}

	return macs
}
