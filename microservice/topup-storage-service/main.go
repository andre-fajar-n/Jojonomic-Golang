package main

import (
	"context"
	"encoding/json"
	"fmt"
	"jojonomic/utils"
	"jojonomic/utils/model"
	"log"
)

func main() {
	utils.InitConfig()

	utils.InitializeDatabase()

	subscribeData()

	fmt.Println("Successful running apps")
}

func subscribeData() {
	reader := utils.GetKafkaReader(utils.Config.Kafka.URL, utils.Config.Kafka.TopicTopup, "")

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
			continue
		}

		err = utils.DB.Create(req).Error
		if err != nil {
			fmt.Println("error insert data to db:", err)
			continue
		}

		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
