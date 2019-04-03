package simple_service

import (
	"context"
	"fmt"
	"github.com/zibilal/go-data-access/persistence"
	"github.com/zibilal/go-data-access/persistence/mongodbpersistence"
)

type FetchService struct {
	dataPersistence persistence.Persistence
}

func NewFetchService(dataPersistence persistence.Persistence) *FetchService {
	fetchService := new(FetchService)
	fetchService.dataPersistence = dataPersistence

	return fetchService
}

func (s *FetchService) Serve(ctx context.Context, data ...interface{}) error {
	if len(data) == 1 {
		fmt.Println("Fetch.Serve")
		qry := make(map[string]interface{})
		err := s.dataPersistence.Fetch(ctx, qry, "find", "simple", data[0])
		if err != nil {
			return err
		}
	} else if len(data) == 2 {
		fmt.Println("Fetch.Serve two inputs")
		qry := mongodbpersistence.QueryMap{}
		err := qry.FromStruct(data[1])
		if err != nil {
			return err
		}

		err = s.dataPersistence.Fetch(ctx, qry, "find", "simple", data[0])
		if err != nil {
			return err
		}
	}
	return nil
}
