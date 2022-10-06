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

	r.HandleFunc("/api/saldo", handlerCheckSaldo).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8005", r))
}

func handlerCheckSaldo(w http.ResponseWriter, r *http.Request) {
	req := model.CheckSaldoRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		fmt.Println("error decode data", err)
		utils.WriteErrorResponse(w, "Bad request", err)
		return
	}
	defer r.Body.Close()

	var data model.TblRekening

	err := utils.DB.Model(&data).Where("norek = ?", req.Norek).First(&data).Error
	if err != nil {
		utils.WriteErrorResponse(w, "", err)
		return
	}

	utils.WriteSuccessResponseWithData(w, map[string]interface{}{
		"norek": data.Norek,
		"saldo": data.GoldBalance,
	}, "")
}
