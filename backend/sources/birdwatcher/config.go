package birdwatcher

type Config struct {
	Id   int
	Name string

	Api            string `ini:"api"`
	Timezone       string `ini:"timezone"`
	ShowLastReboot bool   `ini:"show_last_reboot"`
}
