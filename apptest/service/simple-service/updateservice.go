package simple_service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zibilal/go-data-access/persistence"
	"github.com/zibilal/go-data-access/persistence/mongodbpersistence"
)

type UpdateService struct {
	dataPersistence persistence.Persistence
}

func NewUpdateService(dataPersistence persistence.Persistence) *UpdateService {
	updateService := new(UpdateService)
	updateService.dataPersistence = dataPersistence

	return updateService
}

func (s *UpdateService) Serve(ctx context.Context, data ...interface{}) error {
	if len(data) == 3 {
		fmt.Println("Update data -- ", data[1])
		qry := mongodbpersistence.QueryMap{}
		err := qry.FromStruct(data[2])
		if err != nil {
			return err
		}

		err = s.dataPersistence.Update(ctx, qry, "simple", data[1])
		if err != nil {
			return err
		}

		err = s.dataPersistence.Fetch(ctx, qry, "find", "simple", data[0])
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid service request")
	}

	return nil
}
