package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alice-lg/alice-lg/pkg/backend"
)

var banner = []string{
	"        **        ***               Alice ?VERSION       ",
	"     *****         ***      *                            ",
	"    *  ***          **     ***      Listening on: ?LISTEN",
	"       ***          **      *       Routeservers: ?RSCOUNT",
	"      *  **         **                                   ",
	"      *  **         **    ***        ****       ***      ",
	"     *    **        **     ***      * ***  *   * ***     ",
	"     *    **        **      **     *   ****   *   ***    ",
	"    *      **       **      **    **         **    ***   ",
	"    *********       **      **    **         ********    ",
	"   *        **      **      **    **         *******     ",
	"   *        **      **      **    **         **          ",
	"  *****      **     **      **    ***     *  ****    *   ",
	" *   ****    ** *   *** *   *** *  *******    *******    ",
	"*     **      **     ***     ***    *****      *****     ",
	"*                                                        ",
	" **                                                      ",
}

func printBanner() {
	status, _ := backend.NewAppStatus()
	cfg := backend.AliceConfig
	mapper := strings.NewReplacer(
		"?VERSION", status.Version,
		"?LISTEN", cfg.Server.Listen,
		"?RSCOUNT", strconv.FormatInt(int64(len(cfg.Sources)), 10),
	)

	for _, l := range banner {
		l = mapper.Replace(l)
		fmt.Println(l)
	}
}
