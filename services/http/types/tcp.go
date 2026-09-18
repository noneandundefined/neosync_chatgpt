package types

// SIZE(uint16_t) | TYPE(uint8_t) | REQID(16 bytes) | IMEI(8 bytes) | STATCODE(uint16_t) | DATA(DATA_LEN)

type RabbitMQ_TransitBinary struct {
	Type      uint8
	RequestID string // [16]byte
	Imei      string // [8]byte
	NResp     bool
	StatCode  uint16
	Data      []byte
}
