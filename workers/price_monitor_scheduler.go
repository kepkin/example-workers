package workers

import (
    "context"
    "log"
    "time"

    "github.com/google/uuid"

    "github.com/reactivego/scheduler"
    ipm "github.com/kepkin/interview-price-monitor"
)


func NewPriceMonitorScheduler(repo ipm.Repository, resolver ipm.PriceResolver, concurrent scheduler.ConcurrentScheduler) PriceMonitorScheduler {
    return PriceMonitorScheduler{
        Repo:     repo,
        Resolver: resolver, 

        concurrent:      concurrent,
        runningMonitors: make(map[uuid.UUID]ipm.MonitorTask),
        doneCh:          make(chan uuid.UUID),
    }
}

type PriceMonitorScheduler struct {
    Repo     ipm.Repository
    Resolver ipm.PriceResolver

    concurrent scheduler.ConcurrentScheduler

    runningMonitors map[uuid.UUID]ipm.MonitorTask

    doneCh chan uuid.UUID
}

//TODO: check for panic and reschedule
func (t PriceMonitorScheduler) Do(again func(due time.Duration)) {

    for {
        select {
            case taskID := <-t.doneCh:
                log.Print("removing monitor from runnings ", taskID)
                t.finishTask(taskID)

            case <-time.After(time.Second):
                ;
        }


        t.runNewMonitors()
        t.updateLiveStatus()
    }


    //again(time.Second)
}

func (t PriceMonitorScheduler) runNewMonitors() {
    monitors, err := t.Repo.ListMonitorsToStart(context.Background())
    if err != nil {
        log.Print("failed to get list of tasks ", err)
    }

    for _, v := range monitors {
        if _, ok := t.runningMonitors[v.ID]; ok {
            // it's already running
            continue
        }
        t.runningMonitors[v.ID] = v
        t.schedule(v)
    }   
}

func (t PriceMonitorScheduler) finishTask(id uuid.UUID) {
    delete(t.runningMonitors, id)
    err := t.Repo.UpdateStatus(context.Background(), id, ipm.TaskStatusFinished)
    if err != nil {
        log.Print("update status err ", id, " ", err)
    }
}

func (t PriceMonitorScheduler) updateLiveStatus() {
    for _, v := range t.runningMonitors {
        err := t.Repo.UpdateStatus(context.Background(), v.ID, ipm.TaskStatusRunning)
        if err != nil {
            log.Print("update status err ", v.ID, " ", err)
        }
    }
}

func (t PriceMonitorScheduler) schedule(v ipm.MonitorTask) {
        rt := PriceMonitorWorker{
            priceMonSvc: ipm.PriceMonitorService{
                Repo:     t.Repo,
                Resolver: t.Resolver,
            },
            monitor: v,
            doneCh:  t.doneCh,
        }

        log.Print("scheduling ", v.ID)
        t.concurrent.ScheduleFutureRecursive(v.Start.Sub(time.Now()), rt.Do)

        //TODO: may be this is incorrect place
        err := t.Repo.UpdateStatus(context.Background(), v.ID, ipm.TaskStatusRunning)
        if err != nil {
            log.Print("update status err ", v.ID, " ", err)
        }
}