package elastic

import (
	"context"
	"errors"
	"fmt"
	"product-es-migration/app/product/repository"
	"product-es-migration/app/product/usecase"
	"product-es-migration/database/postgres"
	"product-es-migration/domain"
)

func Migrate() error {
	client, err := connect()
	if err != nil {
		errTemplate := fmt.Sprintf("error connecting to elastic db with details : %v", err)
		return errors.New(errTemplate)
	}
	exist, err := checkExist(client)
	if err != nil {
		return err
	}
	if exist {
		ctx := context.Background()
		dbRepo := postgres.OpenPostgres()

		repositoryCollection := domain.RepositoryCollection{
			ProductRepo: repository.NewProductRepo(dbRepo),
		}

		usecaseCollection := domain.UsecaseCollection{
			ProductUC: usecase.NewProductUC(repositoryCollection),
		}
		page := 1
		for {
			p, err := usecaseCollection.ProductUC.GetProducts(ctx, page, 100)

			if len(p) == 0 || err != nil {
				break
			}
			fmt.Printf("page : %d \n", page)

			bulk(client, p)
			page++
		}
		return nil
	}
	return errors.New("index is not found. please create index before doing migration")
}
