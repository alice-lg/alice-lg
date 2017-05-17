package birdwatcher

type Config struct {
	Api            string `ini:"api"`
	ShowLastReboot bool   `ini:"show_last_reboot"`
}
