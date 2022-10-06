package main

import (
	"encoding/json"
	"fmt"
	"jojonomic/utils"
	"jojonomic/utils/model"
	"log"
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

	r.HandleFunc("/api/input-harga", handlerInputHarga).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8001", r))
}

func handlerInputHarga(w http.ResponseWriter, r *http.Request) {
	var req model.InputHargaRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		fmt.Println("error decode data", err)
		utils.WriteErrorResponse(w, "Bad request", err)
		return
	}
	defer r.Body.Close()

	id, err := shortid.Generate()
	if err != nil {
		fmt.Println("error Generate uuid")
		utils.WriteErrorResponse(w, "", err)
		return
	}

	kafkaWriter := utils.ConnectToKafka(utils.Config.Kafka.URL, utils.Config.Kafka.TopicInputHarga)
	defer kafkaWriter.Close()

	byteData, err := json.Marshal(model.TblHarga{
		ReffID:       id,
		AdminID:      req.AdminID,
		HargaTopup:   req.HargaTopup,
		HargaBuyback: req.HargaBuyback,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		log.Fatal("error marshal:", err)
		utils.WriteErrorResponse(w, id, err)
		return
	}

	msg := kafka.Message{
		Key:   []byte(id),
		Value: byteData,
	}
	_, err = kafkaWriter.WriteMessages(msg)
	if err != nil {
		fmt.Println("error write message in kafka:", err)
		utils.WriteErrorResponse(w, id, err)
		return
	}

	utils.WriteSuccessResponse(w, id)
}
