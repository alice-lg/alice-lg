package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/alice-lg/alice-lg/pkg/config"
	"github.com/alice-lg/alice-lg/pkg/http"
	"github.com/alice-lg/alice-lg/pkg/store"
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

func printBanner(
	cfg *config.Config,
	neighborsStore *store.NeighborsStore,
	routesStore *store.RoutesStore,
) {
	ctx := context.Background()
	status, _ := http.CollectAppStatus(ctx, routesStore, neighborsStore)
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
