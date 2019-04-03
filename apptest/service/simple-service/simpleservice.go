package simple_service

import (
	"context"
	"fmt"
	"github.com/zibilal/go-data-access/persistence"
)

type SimpleService struct {
	dataPersistence persistence.Persistence
}

func NewSimpleService(dataPersistence persistence.Persistence) *SimpleService {
	simpleService := new(SimpleService)
	simpleService.dataPersistence = dataPersistence

	return simpleService
}

func (s *SimpleService) Serve(ctx context.Context, data ...interface{}) error {
	if len(data) > 0 {
		var err error
		for _, d := range data {
			err = s.dataPersistence.Store(ctx, "simple", d)
			if err != nil {
				return err
			}
		}

		fmt.Println("PROCESS SUCCESSFUL")
	} else {
		fmt.Println("NO DATA SUBMITTED")
	}

	return nil
}
