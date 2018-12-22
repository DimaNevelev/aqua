package handler

import (
	"github.com/dimanevelev/travers/model"
	"github.com/dimanevelev/travers/persistence"
	"github.com/dimanevelev/travers/utils"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"sync"
)

type StatsHandler struct {
	persistence.Client
}

func NewStatsHandler(client persistence.Client) *StatsHandler {
	return &StatsHandler{client}
}

func (handler *StatsHandler) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", handler.Get)
	return router
}

func (handler *StatsHandler) Get(w http.ResponseWriter, r *http.Request) {
	stats := model.Stats{}
	var errors [5]error
	var wg sync.WaitGroup

	wg.Add(5)
	go func() {
		stats.TotalFiles, errors[0] = handler.Client.CountRows()
		wg.Done()
	}()
	go func() {
		stats.MaxFile.Size, stats.MaxFile.Path, errors[1] = handler.Client.MaxFileSize()
		wg.Done()
	}()
	go func() {
		stats.AvgFileSize, errors[2] = handler.Client.AVGFileSize()
		wg.Done()
	}()
	go func() {
		stats.Extensions, errors[3] = handler.Client.ExtensionsList()
		wg.Done()
	}()
	go func() {
		stats.TopExtension, errors[4] = handler.Client.MostCommonExt()
		wg.Done()
	}()
	wg.Wait()

	for _, err := range errors {
		if err != nil {
			log.Fatal(err)
			utils.SendHTTPError(w, err, 500)
		}
	}

	utils.SendHTTPResponse(w, stats)
}
