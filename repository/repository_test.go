package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	ipm "github.com/kepkin/interview-price-monitor"
)

const defaultConnString = "postgres://postgres:password@localhost:8432/monitor"


func TestListMonitors(t *testing.T) {
	assert := assert.New(t)

	repo, err := newRepo(defaultConnString)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.resetDB()
	if err != nil {
		t.Fatal(err)
	}

	activeTasks := []ipm.MonitorTask{
		ipm.MonitorTask{
			ID: uuid.New(),
			Start: time.Now().Add(10*time.Second),
			Stop: time.Now().Add(time.Hour),
			Target: "bitcoin",
		},
	}

	nonActiveTasks := []ipm.MonitorTask{
		ipm.MonitorTask{
			ID: uuid.New(),
			Start: time.Now().Add(time.Hour),
			Stop: time.Now().Add(time.Hour*2),
			Target: "bitcoin",
		},
	}

	data := make([]ipm.MonitorTask, len(activeTasks) + len(nonActiveTasks))
	copy(data[0:len(activeTasks)], activeTasks)
	copy(data[len(activeTasks):], nonActiveTasks)

	for _, d := range data {
		err = repo.AddMonitor(context.Background(), d)
		if err != nil {
			t.Fatal(err)
		}
	}

	result, err := repo.ListMonitorsToStart(context.Background()) 
	assert.NoError(err)
	assertElementsMatch(t, activeTasks, result, func(a, b ipm.MonitorTask) bool { return a.ID == b.ID})
}


