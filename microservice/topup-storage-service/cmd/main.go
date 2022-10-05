package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"topup-storage-service/infra"
	"topup-storage-service/model"

	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

func main() {
	r := mux.NewRouter()

	db := infra.InitializeDatabase()

	subscribeData(db)

	err := http.ListenAndServe("localhost:8083", r)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successful running apps")
}

func subscribeData(db *gorm.DB) {
	topic := "topup"

	reader := infra.GetKafkaReader("localhost:9093", topic, "")

	defer reader.Close()

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}

		req := model.TblTransaksi{}
		err = json.Unmarshal(m.Value, &req)
		if err != nil {
			fmt.Println("error unmarshal:", err)
			return
		}

		id, err := generateID()
		if err != nil {
			fmt.Println("error generateID:", err)
			return
		}

		req.ReffID = id
		req.CreatedAt = time.Now().Unix()

		err = db.Create(req).Error
		if err != nil {
			fmt.Println("error insert data to db:", err)
			return
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}

func generateID() (string, error) {
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		fmt.Println("error New shortid")
		return "", err
	}

	id, err := sid.Generate()
	if err != nil {
		fmt.Println("error Generate uuid")
		return "", err
	}

	return id, nil
}
