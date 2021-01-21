package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	contName = flag.StringP("container", "c", "Default container", "The container name")
	interval = flag.IntP("interval", "i", 10, "collection interval (in seconds)")
	output   = flag.StringP("output", "o", "./data/node_status.json", "Output Dir")
)

func main() {
	flag.Parse()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	stats := getStats(cli, *contName)
	start := time.Now()
	startUsage := stats.CPUStats.CPUUsage.TotalUsage
	printStats(stats, *contName, start, startUsage, 1)

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	for range ticker.C {
		stats = getStats(cli, *contName)
		now := time.Now()
		elapsed := now.Sub(start)
		printStats(stats, *contName, now, startUsage, elapsed)
	}
}

func getStats(cli *client.Client, contName string) *types.StatsJSON {
	resp, err := cli.ContainerStats(context.TODO(), contName, false)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	stats := &types.StatsJSON{}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, stats)
	if err != nil {
		panic(err)
	}

	return stats
}

func printStats(stats *types.StatsJSON, contName string, now time.Time, startUsage uint64, elapsed time.Duration) {
	ts := now.UTC().Format(time.RFC3339)
	timeElapsed := elapsed.Seconds()
	nowNetwork := stats.Networks
	log.Info(nowNetwork)
	percentCPUSinceStart := float64(stats.CPUStats.CPUUsage.TotalUsage-startUsage) / float64(elapsed.Nanoseconds()) * 100
	// json
	text := fmt.Sprintf(`{"ts":"%s","c":"%s","timeElapsed":%.2f,"cpu":%.2f,"mUsageMiB":%.1f}`,
		ts,
		contName,
		timeElapsed,
		percentCPUSinceStart,
		float64(stats.MemoryStats.Usage/(1024*1024)))

	err := ioutil.WriteFile(*output, []byte(text), 0777)
	if err != nil {
		log.Info(err)
	}
}
