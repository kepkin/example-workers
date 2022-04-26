package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	ipm "github.com/kepkin/interview-price-monitor"
)

func (s *PriceMonitorServiceServerImpl) GetMonitor(in GetMonitorRequest, c *gin.Context) {
	id, err := uuid.Parse(in.Path.MonitorID)
	if err != nil {
		//TODO: add meaningfull error
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := s.repo.ListMonitorResults(context.TODO(), id)
	c.JSON(http.StatusOK, res)
}

func (s *PriceMonitorServiceServerImpl) StartMonitor(in StartMonitorRequest, c *gin.Context) {
	task := in.Body.JSON

	id, err := uuid.Parse(in.Path.MonitorID)
	if err != nil {
		//TODO: add meaningfull error
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	frequency, err := time.ParseDuration(task.Frequency)
	if err != nil {
		//TODO: add meaningfull error
		c.AbortWithStatus(http.StatusBadRequest)
		return		
	}

	err = s.repo.AddMonitor(context.TODO(), ipm.MonitorTask{
		ID: id,
		Start: task.Start,
		Stop: task.Stop,
		Target: task.Target,
		Frequency: frequency,
	})
	if err != nil {
		//TODO: add meaningfull error
		log.Print(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
