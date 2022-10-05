package main

import (
	"encoding/json"
	"fmt"
	"jojonomic/utils"
	"log"
	"net/http"
	"topup-service/infra"
	"topup-service/model"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

func main() {
	r := mux.NewRouter()

	infra.InitializeDatabase()

	r.HandleFunc("/api/topup", handlerTopup).Methods("POST")

	err := http.ListenAndServe("localhost:8082", r)
	if err != nil {
		panic(err)
	}
}

func handlerTopup(w http.ResponseWriter, r *http.Request) {
	var req model.TopupRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		fmt.Println("error decode data", err)
		utils.WriteErrorResponse(w, "Bad request", err)
		return
	}
	defer r.Body.Close()

	var harga model.TblHarga
	err := infra.DB.Model(&model.TblHarga{}).Last(&harga).Error
	if err != nil {
		fmt.Println("error get latest harga:", err)
		utils.WriteErrorResponse(w, "Bad request", err)
		return
	}

	if harga.HargaTopup != req.Harga {
		err = fmt.Errorf("harga topup tidak sesuai dengan harga topup saat ini")
		utils.WriteErrorResponse(w, "Bad request", err)
		return
	}

	var rekening model.TblRekening
	err = infra.DB.Find(&rekening).Where("norek = $1", req.Norek).Error
	if err != nil {
		utils.WriteErrorResponse(w, "Bad request", err)
		return
	}

	kafkaWriter := infra.GetKafkaWriter("topup")
	defer kafkaWriter.Close()

	byteData, err := json.Marshal(model.TblTransaksi{
		Norek:        req.Norek,
		Type:         "topup",
		GoldWeight:   req.Gram,
		HargaTopup:   harga.HargaTopup,
		HargaBuyback: harga.HargaBuyback,
		GoldBalance:  rekening.GoldBalance + req.Gram,
	})
	if err != nil {
		log.Fatal("error marshal:", err)
		utils.WriteErrorResponse(w, "Bad request", err)
		return
	}

	msg := kafka.Message{
		Key:   []byte("topup"),
		Value: byteData,
	}
	err = kafkaWriter.WriteMessages(r.Context(), msg)

	if err != nil {
		utils.WriteErrorResponse(w, "Bad request", err)
		w.Write([]byte(err.Error()))
		log.Fatalln(err)
	}

	utils.WriteSuccessResponse(w, "Success topup")
}
