package backend

import (
	"fmt"
	"strconv"
	"strings"
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
	status, _ := NewAppStatus()
	mapper := strings.NewReplacer(
		"?VERSION", status.Version,
		"?LISTEN", AliceConfig.Server.Listen,
		"?RSCOUNT", strconv.FormatInt(int64(len(AliceConfig.Sources)), 10),
	)

	for _, l := range banner {
		l = mapper.Replace(l)
		fmt.Println(l)
	}
}
