package birdwatcher

type Config struct {
	Id   string
	Name string

	Api             string `ini:"api"`
	Timezone        string `ini:"timezone"`
	ServerTime      string `ini:"servertime"`
	ServerTimeShort string `ini:"servertime_short"`
	ServerTimeExt   string `ini:"servertime_ext"`
	ShowLastReboot  bool   `ini:"show_last_reboot"`

	Type               string `ini:"type"`
	PeerTablePrefix    string `ini:"peer_table_prefix"`
	PipeProtocolPrefix string `ini:"pipe_protocol_prefix"`
}
