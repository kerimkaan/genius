package types

// NTPConfiguration represents a single NTP server entry from the configuration file.
type NTPConfiguration struct {
	Server string `json:"server"`
	IBurst bool   `json:"iburst"`
}

// FolderSize holds a folder path and its total size in bytes.
type FolderSize struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}
