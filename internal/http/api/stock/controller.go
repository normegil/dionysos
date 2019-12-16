package stock

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/normegil/dionysos"
	"log"
	"net/http"
	"strconv"
)

func NewController() http.Handler {
	rt := chi.NewRouter()
	rt.Get("/", LoadAll)
	return rt
}

func LoadAll(w http.ResponseWriter, r *http.Request) {
	items := make([]dionysos.Item, 0)
	for i := 0; i < 10; i++ {
		items = append(items, dionysos.Item{
			ID:   uuid.New(),
			Name: "Item" + strconv.Itoa(i),
		})
	}

	bytes, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); nil != err {
		log.Print(err)
		return
	}
}
