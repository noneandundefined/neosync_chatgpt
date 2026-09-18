package rabbitmq

import (
	"encoding/binary"
	"fmt"
	"neomatica/neosync/types"

	"github.com/google/uuid"
)

const (
	REQID_SIZE = 16
	IMEI_SIZE  = 8
)

func uuidToBytes(uuidStr string) ([REQID_SIZE]byte, error) {
	u, err := uuid.Parse(uuidStr)
	if err != nil {
		return [REQID_SIZE]byte{}, err
	}

	var out [REQID_SIZE]byte
	copy(out[:], u[:])
	return out, nil
}

func charToNibble(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	return 0xF
}

func imeiToBytes(imei string) ([IMEI_SIZE]byte, error) {
	var out [IMEI_SIZE]byte

	if len(imei) > 16 {
		return out, fmt.Errorf("imei too long")
	}

	padded := imei
	if len(padded)%2 != 0 {
		padded += "F"
	}

	for i := 0; i < len(padded)/2 && i < IMEI_SIZE; i++ {
		hi := padded[i*2]
		lo := padded[i*2+1]

		out[i] = (charToNibble(hi) << 4) | charToNibble(lo)
	}

	return out, nil
}

func bcdToIMEI(b [IMEI_SIZE]byte) string {
	out := make([]byte, 0, 16)

	for _, v := range b {
		hi := v >> 4
		lo := v & 0x0F

		out = append(out, nibbleToChar(hi))
		if lo != 0xF {
			out = append(out, nibbleToChar(lo))
		}
	}

	return string(out)
}

func nibbleToChar(n byte) byte {
	return '0' + n
}

func (r *RabbitMQ) Encoder(rabbitmqTransit types.RabbitMQ_TransitBinary) ([]byte, error) {
	/* RequestID */
	reqIdBytes, err := uuidToBytes(rabbitmqTransit.RequestID)
	if err != nil {
		return nil, err
	}

	/* Imei */
	imeiBytes, err := imeiToBytes(rabbitmqTransit.Imei)
	if err != nil {
		return nil, err
	}

	/* Size packet */
	sizePacket := 1 + REQID_SIZE + IMEI_SIZE + 1 + 2 + len(rabbitmqTransit.Data)

	/* Buffer */
	buffer := make([]byte, 0, sizePacket+2)
	tmpBuffer := make([]byte, 2)

	/* Size write */
	binary.LittleEndian.PutUint16(tmpBuffer, uint16(sizePacket))
	buffer = append(buffer, tmpBuffer...)

	/* Type write */
	buffer = append(buffer, rabbitmqTransit.Type)

	/* RequestID write */
	buffer = append(buffer, reqIdBytes[:]...)

	/* Imei write */
	buffer = append(buffer, imeiBytes[:]...)

	if rabbitmqTransit.NResp {
		buffer = append(buffer, 1)
	} else {
		buffer = append(buffer, 0)
	}

	/* StatCode write */
	binary.LittleEndian.PutUint16(tmpBuffer, rabbitmqTransit.StatCode)
	buffer = append(buffer, tmpBuffer...)

	/* Data write */
	buffer = append(buffer, rabbitmqTransit.Data...)

	return buffer, nil
}

func (r *RabbitMQ) Decoder(buffer []byte) (types.RabbitMQ_TransitBinary, error) {
	offset := 0

	if len(buffer) < 2 {
		return types.RabbitMQ_TransitBinary{}, fmt.Errorf("buffer too short to read size")
	}

	/* Size read */
	sizePacket := binary.LittleEndian.Uint16(buffer[offset:])
	offset += 2

	if int(sizePacket) != len(buffer)-2 {
		return types.RabbitMQ_TransitBinary{}, fmt.Errorf("size mismatch: expected %d, got %d", sizePacket, len(buffer)-2)
	}

	/* Type read */
	type_p := buffer[offset]
	offset++

	/* RequestId read */
	if len(buffer) < offset+REQID_SIZE {
		return types.RabbitMQ_TransitBinary{}, fmt.Errorf("buffer too short to read RequestID")
	}

	var requestId [REQID_SIZE]byte
	copy(requestId[:], buffer[offset:offset+REQID_SIZE])
	offset += REQID_SIZE

	/* Imei read */
	if len(buffer) < offset+IMEI_SIZE {
		return types.RabbitMQ_TransitBinary{}, fmt.Errorf("buffer too short to read IMEI")
	}

	var imei [IMEI_SIZE]byte
	copy(imei[:], buffer[offset:offset+IMEI_SIZE])
	offset += IMEI_SIZE

	/* Need response */
	if len(buffer) < offset+1 {
		return types.RabbitMQ_TransitBinary{}, fmt.Errorf("buffer too short to read NResp")
	}
	nresp := buffer[offset] == 1
	offset++

	/* StatCode read */
	if len(buffer) < offset+2 {
		return types.RabbitMQ_TransitBinary{}, fmt.Errorf("buffer too short to read StatCode")
	}

	statCode := binary.LittleEndian.Uint16(buffer[offset:])
	offset += 2

	/* Data read */
	var data []byte
	if offset < len(buffer) {
		data = buffer[offset:]
	} else {
		data = []byte{}
	}

	return types.RabbitMQ_TransitBinary{
		Type:      type_p,
		RequestID: uuid.UUID(requestId).String(),
		Imei:      bcdToIMEI(imei),
		NResp:     nresp,
		StatCode:  statCode,
		Data:      data,
	}, nil
}
