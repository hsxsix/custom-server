package main

import (
	"custom_server/internal/api"
	"custom_server/pkg/server"
	"flag"
	"fmt"
	"strings"
)

var (
	// build cmd: go build -a -ldflags "-X 'main.gitBranch=$(git rev-parse --abbrev-ref HEAD)' -X 'main.gitCommit=$(git show -s --format=%s)' -X 'main.goVersion=$(go version)' -X 'main.buildTime=$(date '+%Y-%m-%d %H:%M:%S')' -extldflags -static" -o web_server
	buildTime string
	goVersion string
	gitBranch string
	gitHash   string
)

func main() {
	var (
		debug   bool
		version bool
		config  string
	)

	flag.BoolVar(&version, "v", false, "print version info")
	flag.BoolVar(&version, "version", false, "print version info")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.StringVar(&config, "c", "config/config.toml", "config file path")
	flag.Parse()

	if version {
		fmt.Printf("******Server Version Info******\n")
		fmt.Printf("Git Branch:%s \n", gitBranch)
		fmt.Printf("Last Commit Hash:%s \n", gitHash)
		fmt.Printf("Last Build Time: %s \n", buildTime)
		fmt.Printf("Go Version: %s\n", strings.Replace(goVersion, "go version ", "", 1))
		return
	}

	web := server.NewWebServer(config, debug).WithName("custom_server")

	// it will initialize the database which named postgres
	web.WithDB("postgres")

	//it will initialize the cron
	//web.WithCron()

	// register router
	web.RegisterRouter(api.ServerRouter)

	web.Run()
}
