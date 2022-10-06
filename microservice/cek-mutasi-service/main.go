package main

import (
	"encoding/json"
	"fmt"
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

	r.HandleFunc("/api/mutasi", handlerCheckSaldo).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8006", r))
}

func handlerCheckSaldo(w http.ResponseWriter, r *http.Request) {
	req := model.CheckMutasiRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		fmt.Println("error decode data", err)
		utils.WriteErrorResponse(w, "Bad request", err)
		return
	}
	defer r.Body.Close()

	var data []model.CheckMutasiResponse

	err := utils.DB.Model(&model.TblTransaksi{}).
		Where("norek = ?", req.Norek).
		Where("created_at >= ?", req.StartDate).
		Where("created_at <= ?", req.EndDate).
		Order("created_at DESC").
		Find(&data).Error
	if err != nil {
		utils.WriteErrorResponse(w, "", err)
		return
	}

	utils.WriteSuccessResponseWithData(w, data, "")
}
