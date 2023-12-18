package elastic

import (
	"errors"
	"fmt"
	"os"
	"product-es-migration/config"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

const NotFound = 404

func checkExist(client *elasticsearch.Client) (bool, error) {
	response, err := client.Indices.Exists([]string{config.ELASTICSEARCH_PRODUCT_INDEX})
	if err != nil {
		errTemplate := fmt.Sprintf("error checking index with details : %v", err)
		return false, errors.New(errTemplate)
	}
	if response.IsError() {
		if response.StatusCode == NotFound {
			return false, nil
		} else {
			errTemplate := fmt.Sprintf("error checking index with details : %v", response.String())
			return false, errors.New(errTemplate)
		}
	}
	return true, nil
}

func deleteIndex(client *elasticsearch.Client) error {
	exist, err := checkExist(client)
	if err != nil {
		return err
	}
	if exist {
		response, err := client.Indices.Delete([]string{config.ELASTICSEARCH_PRODUCT_INDEX})
		if err != nil {
			return err
		}
		if response.IsError() {
			if response.StatusCode == NotFound {
				return errors.New("error deleting index with error : index is not found")
			} else {
				errTemplate := fmt.Sprintf("error deleting index with details : %v", response.String())
				return errors.New(errTemplate)
			}
		}

		fmt.Printf("Index %v is deleted ", config.ELASTICSEARCH_PRODUCT_INDEX)
		return nil
	}
	errTemplate := fmt.Sprintf("index %v is not exist", config.ELASTICSEARCH_PRODUCT_INDEX)
	return errors.New(errTemplate)

}

func createIndex(client *elasticsearch.Client) error {
	mapper, err := os.Open("mapper.yaml")
	if err != nil {
		errTemplate := fmt.Sprintf("error getting schema in mapper.yaml with error : %v", err)
		return errors.New(errTemplate)
	}
	defer mapper.Close()

	exist, err := checkExist(client)
	if err != nil {
		return err
	}
	if !exist {
		res, err := client.Indices.Create(
			config.ELASTICSEARCH_PRODUCT_INDEX,
			client.Indices.Create.WithBody(mapper),
		)
		if err != nil {
			return err
		}
		if res.IsError() {
			return errors.New(res.String())
		}

		fmt.Printf("Index %v is created ", config.ELASTICSEARCH_PRODUCT_INDEX)
		return nil
	}

	errTemplate := fmt.Sprintf("index is exist with name:  %v", config.ELASTICSEARCH_PRODUCT_INDEX)
	return errors.New(errTemplate)
}

func CreateIndex() error {
	client, err := connect()
	if err != nil {
		errTemplate := fmt.Sprintf("error connecting to elastic db with details : %v", err)
		return errors.New(errTemplate)
	}
	mapper, err := os.Open("mapper.yaml")
	if err != nil {
		errTemplate := fmt.Sprintf("elastic mapper schema not found : %v", err)
		return errors.New(errTemplate)
	}
	defer mapper.Close()

	return createIndex(client)
}

func DeleteIndex() error {
	client, err := connect()
	if err != nil {
		errTemplate := fmt.Sprintf("error connecting to elastic db with details : %v", err)
		return errors.New(errTemplate)
	}

	return deleteIndex(client)
}

func CheckIndex() error {
	client, err := connect()
	if err != nil {
		errTemplate := fmt.Sprintf("error connecting to elastic db with details : %v", err)
		return errors.New(errTemplate)
	}

	any, err := checkExist(client)
	if err != nil {
		errTemplate := fmt.Sprintf("error checking index with details : %v", err)
		return errors.New(errTemplate)
	}
	if any {
		msgTemplate := fmt.Sprintf("index %v is found", config.ELASTICSEARCH_PRODUCT_INDEX)
		fmt.Print(msgTemplate)
		return nil
	} else {
		msgTemplate := fmt.Sprintf("index %v is not found", config.ELASTICSEARCH_PRODUCT_INDEX)
		fmt.Print(msgTemplate)
		return nil
	}
}
