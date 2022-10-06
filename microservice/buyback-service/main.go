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

	kafkaWriter := utils.ConnectToKafka(utils.Config.Kafka.URL, utils.Config.Kafka.TopicBuyback)
	defer kafkaWriter.Close()

	r.HandleFunc("/api/buyback", handlerTopup(kafkaWriter)).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8003", r))
}

func handlerTopup(kafkaWriter *kafka.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
			utils.WriteErrorResponse(w, "", fmt.Errorf("buyback harus kelipatan 0.001"))
			return
		}

		// validate harga from table and request
		if harga.HargaBuyback != req.Harga {
			err = fmt.Errorf("harga buyback tidak sesuai dengan harga buyback saat ini")
			utils.WriteErrorResponse(w, "", err)
			return
		}

		// get data rekening
		var rekening model.TblRekening
		err = utils.DB.Find(&rekening).Where("norek = $1", req.Norek).Error
		if err != nil {
			utils.WriteErrorResponse(w, "", err)
			return
		}

		if rekening.GoldBalance < req.Gram {
			utils.WriteErrorResponse(w, "", fmt.Errorf("saldo emas tidak cukup"))
			return
		}

		id, err := shortid.Generate()
		if err != nil {
			fmt.Println("error Generate uuid")
			utils.WriteErrorResponse(w, "", err)
			return
		}

		byteData, err := json.Marshal(model.TblTransaksi{
			ReffID:       id,
			Norek:        req.Norek,
			Type:         "buyback",
			GoldWeight:   req.Gram,
			HargaTopup:   harga.HargaTopup,
			HargaBuyback: harga.HargaBuyback,
			GoldBalance:  rekening.GoldBalance - req.Gram,
			CreatedAt:    time.Now().Unix(),
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
			utils.WriteErrorResponse(w, id, err)
			return
		}

		utils.WriteSuccessResponse(w, id)
	}
}
