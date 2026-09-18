package adm

type Adm struct{}

type FieldConfiguration struct {
	UID   uint16 `json:"uid" validate:"required,gt=0"`
	Size  uint8  `json:"size" validate:"required,gt=0"`
	Value []byte `json:"value" validate:"required"`
}

type FieldConfigurationParsed struct {
	UID      string `json:"uid" validate:"required,gt=0"`
	RawUID   uint16 `json:"raw_uid"`
	Size     uint8  `json:"size" validate:"required,gt=0"`
	Value    any    `json:"value" validate:"required"`
	RawValue []byte `json:"raw_value"`
	Unknown  bool   `json:"unknown"`
}

type FieldSchema struct {
	UID        uint16 `json:"uid"`
	ResolveUID uint16 `json:"resolve_uid"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Min        *int   `json:"min,omitempty"`
	Max        *int   `json:"max,omitempty"`
	MinLength  *int   `json:"min_length,omitempty"`
	MaxLength  *int   `json:"max_length,omitempty"`
	MaxItems   *int   `json:"max_items,omitempty"`
	IsDigits   *bool  `json:"is_digits,omitempty"`
	Section    string `json:"section"`
}

type PacketConfiguration struct {
	Magic            uint16 `json:"magic"`
	LastModification string `json:"lastModification"`
	Version          uint16 `json:"version"`
	CfgSize          uint16 `json:"cfgSize"`
	CfgHash          uint32 `json:"cfgHash"`
	PacketFieldConfiguration
}

type PacketFieldConfiguration struct {
	UID      uint16 `json:"uid"`
	RawUID   uint16 `json:"raw_uid"`
	Size     uint8  `json:"size"`
	Value    any    `json:"value"`
	RawValue []byte `json:"raw_value"`
	Unknown  bool   `json:"unknown"`
}
