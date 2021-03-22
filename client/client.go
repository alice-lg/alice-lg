package client

import (
	"embed"
)

// Assets hold the alice-lg frontend build
//go:embed build/*
var Assets embed.FS
