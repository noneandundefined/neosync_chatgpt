package adm_v1

func (adm *ADM_V1) Encode_CommandPacket(buffer []byte) []byte {
	if buffer == nil {
		return nil
	}

	sendBuffer := append(buffer, 0x00)

	return sendBuffer
}

func (adm *ADM_V1) Encode_GetSyncPacket() []byte {
	return nil
}

func (adm *ADM_V1) Encode_GetCfgPacket() []byte {
	return nil
}

func (adm *ADM_V1) Encode_SetCfgPacket(buffer []byte, imei *string, hash uint32) []byte {
	return nil
}
