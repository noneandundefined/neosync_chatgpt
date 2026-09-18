package types

/* Структура для парсинга пакета конфигурации */
type ParseConfiguration struct {
	Magic            uint16 `json:"magic"`
	LastModification string `json:"lastModification"`
	Version          uint16 `json:"version"`
	CfgSize          uint16 `json:"cfgSize"`
	CfgHash          uint32 `json:"cfgHash"`
	ParsedConfiguration
}

/* Структура для парсинга пакета конфигурации */
type ParsedConfiguration struct {
	UID   uint16 `json:"uid"`
	Size  uint8  `json:"size"`
	Value any    `json:"value"`
}
