package handler

import (
	"encoding/json"
	"errors"
	"github.com/dimanevelev/travers/model"
	"github.com/dimanevelev/travers/persistence"
	"github.com/dimanevelev/travers/utils"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

type FileHandler struct {
	persistence.Client
}

func NewFileHandler(client persistence.Client) *FileHandler {
	return &FileHandler{client}
}

func (handler *FileHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", handler.Add)
	return router
}

func (handler *FileHandler) Add(w http.ResponseWriter, r *http.Request) {
	file := model.File{}
	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		log.Println("Error occurred while parsing request", err.Error())
		utils.SendHTTPError(w, err, 500)
		return
	}

	err = validateInput(file)
	if err != nil {
		log.Println("Error occurred while validating request", err.Error())
		utils.SendHTTPError(w, err, 400)
		return
	}

	err = handler.Client.Store(file)
	if err != nil {
		log.Println("Error occurred while storing file info", err.Error())
		utils.SendHTTPError(w, err, 500)
		return
	}
	utils.SendHTTPResponse(w, "")
}

func validateInput(file model.File) error {
	if file.Path == "" {
		return errors.New("path can't be empty")
	}
	if file.FileInfo.Name == "" {
		return errors.New("malformed fileInfo. Name can't be empty")
	}
	return nil
}
