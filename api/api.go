package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	ipm "github.com/kepkin/interview-price-monitor"
)

//go:generate go run ./gen/main.go

var _ PriceMonitorService = (*PriceMonitorServiceServerImpl)(nil)

type PriceMonitorServiceServerImpl struct {
	repo ipm.Repository
}

func NewPriceMonitorServiceServer(repo ipm.Repository) PriceMonitorService {
	return &PriceMonitorServiceServerImpl{
		repo: repo,
	}
}

// Errors processing

func (p *PriceMonitorServiceServerImpl) ProcessMakeRequestErrors(c *gin.Context, errors []FieldError) {
	c.JSON(http.StatusBadRequest, fmt.Sprintf("parse request error: %+v", errors))
}

func (p *PriceMonitorServiceServerImpl) ProcessValidateErrors(c *gin.Context, errors []FieldError) {
	c.JSON(http.StatusBadRequest, fmt.Sprintf("validate request error: %+v", errors))
}
