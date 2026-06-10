package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type PaymentServiceSuite struct {
	suite.Suite
	ctx     context.Context
	service *PaymentService
}

func (p *PaymentServiceSuite) SetupTest() {
	p.ctx = context.Background()
	p.service = NewPaymentService()
}

func (p *PaymentService) TearDownTest() {
}

func TestService(t *testing.T) {
	suite.Run(t, new(PaymentServiceSuite))
}
