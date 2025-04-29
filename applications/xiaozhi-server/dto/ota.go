package dto

type OtaResponse struct {
	ServerTime ServerTime  `json:"server_time"`
	Activation *Activation `json:"activation,omitempty"`
	Firmware   Firmware    `json:"firmware"`
	WebSocket  WebSocket   `json:"websocket"`
}

type ServerTime struct {
	Timestamp      int64 `json:"timestamp"`
	TimezoneOffset int32 `json:"timezone_offset"`
}

type Activation struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Challenge string `json:"challenge"`
}

type Firmware struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

type WebSocket struct {
	URL string `json:"url"`
}

type OtaRequest struct {
	Version             int         `json:"version"`
	Language            string      `json:"language"`
	FlashSize           int         `json:"flash_size"`
	MinimumFreeHeapSize int         `json:"minimum_free_heap_size"`
	MacAddress          string      `json:"mac_address"`
	UUID                string      `json:"uuid"`
	ChipModelName       string      `json:"chip_model_name"`
	ChipInfo            ChipInfo    `json:"chip_info"`
	Application         Application `json:"application"`
	PartitionTable      []Partition `json:"partition_table"`
	Ota                 Ota         `json:"ota"`
	Board               Board       `json:"board"`
}

type ChipInfo struct {
	Model    int `json:"model"`
	Cores    int `json:"cores"`
	Revision int `json:"revision"`
	Features int `json:"features"`
}

type Application struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	CompileTime string `json:"compile_time"`
	IdfVersion  string `json:"idf_version"`
	ElfSha256   string `json:"elf_sha256"`
}

type Partition struct {
	Label   string `json:"label"`
	Type    int    `json:"type"`
	Subtype int    `json:"subtype"`
	Address int    `json:"address"`
	Size    int    `json:"size"`
}

type Ota struct {
	Label string `json:"label"`
}

type Board struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	SSID    string `json:"ssid"`
	RSSI    int    `json:"rssi"`
	Channel int    `json:"channel"`
	IP      string `json:"ip"`
	Mac     string `json:"mac"`
}
