package main

import (
	"jojonomic/utils"
	"jojonomic/utils/model"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	utils.InitConfig()

	utils.InitializeDatabase()

	r.HandleFunc("/api/check-harga", handlerCheckHarga).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8004", r))
}

func handlerCheckHarga(w http.ResponseWriter, r *http.Request) {
	var data model.TblHarga

	err := utils.DB.Model(&data).Order("created_at DESC").First(&data).Error
	if err != nil {
		utils.WriteErrorResponse(w, "", err)
		return
	}

	utils.WriteSuccessResponseWithData(w, map[string]interface{}{
		"harga_buyback": data.HargaBuyback,
		"harga_topup":   data.HargaTopup,
	}, "")
}
