package catalogcore

import (
	"monstrolingo_backend/catalog"
)

type Service = catalog.Service

func GetService() (*Service, error) {
	return catalog.GetService()
}
