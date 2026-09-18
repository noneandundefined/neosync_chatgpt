package types

type AdmLogItem struct {
	ID      string `json:"id"`
	TS      int64  `json:"ts"`
	Event   string `json:"event"`
	Payload any    `json:"payload"`
}

type AdmLogsResponse struct {
	Items []AdmLogItem `json:"items"`
}
