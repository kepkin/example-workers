package ipm

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shopspring/decimal"
)

type TaskStatus uint16

const (
	TaskStatusFree TaskStatus = iota
	TaskStatusRunning
	TaskStatusFinished
)

type MonitorTask struct {
	ID        uuid.UUID
	Start     time.Time
	Stop      time.Time
	Frequency time.Duration
	Target    string
	Status    TaskStatus
}

type Decimal decimal.Decimal

func (d *Decimal) String() string {
	return (*decimal.Decimal)(d).String()
}

func (p *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	return (*decimal.Decimal)(p).UnmarshalJSON(decimalBytes)
}

func (p Decimal) MarshalJSON() ([]byte, error) {
	return decimal.Decimal(p).MarshalJSON()
}

func (p *Decimal) Scan(value interface{}) error {
	return (*decimal.Decimal)(p).Scan(value)
}

type PriceMap map[time.Time]Decimal

type PriceResolver interface {
	Get(target string) (time.Time, Decimal, error)
}

type PriceMonitorService struct {
	Resolver PriceResolver
	Repo     Repository
}

func (s *PriceMonitorService) Do(task MonitorTask) error {
	time, price, err := s.Resolver.Get(task.Target)
	if err != nil {
		return err
	}

	return s.Repo.Put(context.TODO(), task.ID, time, price)
}

type Repository interface {
	AddMonitor(ctx context.Context, m MonitorTask) error
	Put(ctx context.Context, id uuid.UUID, time time.Time, price Decimal) error
	ListMonitorsToStart(ctx context.Context) ([]MonitorTask, error)
	ListMonitorResults(ctx context.Context, id uuid.UUID) (PriceMap, error)

	//TODO: this functionality might be splitted to business Service
	CheckLiveProbe(ctx context.Context) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status TaskStatus) error
}
