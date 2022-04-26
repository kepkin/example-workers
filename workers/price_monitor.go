package workers

import (
	"log"
	"time"

	"github.com/google/uuid"

	ipm "github.com/kepkin/interview-price-monitor"
)



type PriceMonitorWorker struct {
	priceMonSvc ipm.PriceMonitorService
	monitor     ipm.MonitorTask

	doneCh chan <- uuid.UUID
}

func (t PriceMonitorWorker) Do(again func(due time.Duration)) {
	err := t.priceMonSvc.Do(t.monitor)
	if err != nil {
		log.Print("task ", t.monitor.ID, " failed to get monitor price ", err)
	}

	if t.monitor.Stop.After(time.Now()) {
		again(t.monitor.Frequency)
	} else {
		t.doneCh <- t.monitor.ID
	}
}
