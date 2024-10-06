package store

import (
	"encoding/json"
	"os"

	"github.com/bebrws/goPR/internal/models"
	"github.com/sirupsen/logrus"
)

func WriteState(stateFilePath string, state *models.GHState) error {
	stateData, err := json.Marshal(state)
	if err != nil {
		return err
	}
	err = os.WriteFile(stateFilePath, stateData, 0644)
	if err != nil {
		return err
	}
	logrus.Info("State written to file: ", stateFilePath)
	return nil
}
