package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/gin-gonic/gin"
	"github.com/reactivego/scheduler"

	ipm "github.com/kepkin/interview-price-monitor"
	"github.com/kepkin/interview-price-monitor/api"
	"github.com/kepkin/interview-price-monitor/repository"
	"github.com/kepkin/interview-price-monitor/workers"
)

func main() {
	var args struct {
		Port int `default:"8080"`
		DB   string `default:"postgres://postgres:password@localhost:8432/monitor"`
	}
	arg.MustParse(&args)

	repo, err := repository.NewRepo(args.DB)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	apiImpl := api.NewPriceMonitorServiceServer(repo)
	api.RegisterRoutes(r, apiImpl)

	startWorkers(repo)

	r.Run(fmt.Sprintf(":%v", args.Port))
}

func startWorkers(repo ipm.Repository) {

	concurrent := scheduler.Goroutine

	priceMonitorScheduler := workers.NewPriceMonitorScheduler(repo, ipm.MockPriceResolver{}, concurrent)

	checkLiveProbe := workers.CheckProbeWorker{
		Repo: repo,
	}

	concurrent.ScheduleFutureRecursive(time.Second, priceMonitorScheduler.Do)
	concurrent.ScheduleFutureRecursive(time.Second, checkLiveProbe.Do)
}
