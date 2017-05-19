package birdwatcher

type Config struct {
	Api            string `ini:"api"`
	Timezone       string `ini:"timezone"`
	ShowLastReboot bool   `ini:"show_last_reboot"`
}
