package main

import (
	"encoding/json"
	"fmt"
	"jojonomic/utils"
	"jojonomic/utils/model"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"github.com/teris-io/shortid"
)

func main() {
	r := mux.NewRouter()

	utils.InitConfig()

	utils.InitializeDatabase()

	r.HandleFunc("/api/topup", handlerTopup).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8002", r))
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

	// get data harga
	var harga model.TblHarga
	err := utils.DB.Model(&model.TblHarga{}).Order("created_at DESC").Find(&harga).Error
	if err != nil {
		fmt.Println("error get latest harga:", err)
		utils.WriteErrorResponse(w, "", err)
		return
	}

	temp := 1000 * req.Gram
	if temp != math.Trunc(temp) {
		utils.WriteErrorResponse(w, "", fmt.Errorf("topup harus kelipatan 0.001"))
		return
	}

	// validate harga from table and request
	if harga.HargaTopup != req.Harga {
		err = fmt.Errorf("harga topup tidak sesuai dengan harga topup saat ini")
		utils.WriteErrorResponse(w, "", err)
		return
	}

	// get data rekening
	var rekening model.TblRekening
	err = utils.DB.Where("norek = $1", req.Norek).Find(&rekening).Error
	if err != nil {
		utils.WriteErrorResponse(w, "", err)
		return
	}
	fmt.Println("REKENING", rekening, req.Norek)

	id, err := shortid.Generate()
	if err != nil {
		fmt.Println("error Generate uuid")
		utils.WriteErrorResponse(w, "", err)
		return
	}

	kafkaWriter := utils.GetKafkaWriter(utils.Config.Kafka.URL, utils.Config.Kafka.TopicTopup)
	defer kafkaWriter.Close()

	byteData, err := json.Marshal(model.TblTransaksi{
		ReffID:       id,
		Norek:        req.Norek,
		Type:         "topup",
		GoldWeight:   req.Gram,
		HargaTopup:   harga.HargaTopup,
		HargaBuyback: harga.HargaBuyback,
		GoldBalance:  rekening.GoldBalance + req.Gram,
		CreatedAt:    time.Now().Unix(),
	})
	fmt.Println("CEK", rekening.GoldBalance, req.Gram, rekening.GoldBalance+req.Gram)
	if err != nil {
		log.Fatal("error marshal:", err)
		utils.WriteErrorResponse(w, id, err)
		return
	}

	msg := kafka.Message{
		Key:   []byte(id),
		Value: byteData,
	}
	err = kafkaWriter.WriteMessages(r.Context(), msg)
	if err != nil {
		utils.WriteErrorResponse(w, id, err)
		return
	}

	utils.WriteSuccessResponse(w, id)
}
