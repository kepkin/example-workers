package workers

import (
	"context"
	"log"
	"time"

	ipm "github.com/kepkin/interview-price-monitor"

)

type CheckProbeWorker struct {
	Repo ipm.Repository
}

func (t CheckProbeWorker) Do(again func(due time.Duration)) {
	defer func() {
		again(time.Second)
	}()

	n, err := t.Repo.CheckLiveProbe(context.Background())
	if err != nil {
		log.Print("CheckLiveProbe failed ", err)
	}
	if n > 0 {
		log.Print(n, " tasks were freed")
	}
}
