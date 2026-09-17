package meta_handler_v1

type ReleaseFile struct {
	Releases []Release `yaml:"releases" json:"releases"`
}

type Release struct {
	Version string                   `yaml:"version" json:"version"`
	Date    string                   `yaml:"date" json:"date"`
	Changes map[string]ChangeSection `yaml:"changes" json:"changes"`
}

type ChangeSection struct {
	Added []string `yaml:"added" json:"added"`
	Fixed []string `yaml:"fixed" json:"fixed"`
}
