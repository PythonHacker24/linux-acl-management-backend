package utils

import (
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/models"
)

/* yaml file loader for config */
func LoadConfig(filename string) (*models.Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var config models.Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}

/* generate a new uuid */
func GenerateTxnID() string {
	return uuid.New().String()
}
