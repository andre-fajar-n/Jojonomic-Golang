module topup-service

go 1.19

require (
	github.com/gorilla/mux v1.8.0
	github.com/segmentio/kafka-go v0.4.35
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569
	jojonomic/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.12.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/klauspost/compress v1.15.7 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/text v0.3.7 // indirect
	gorm.io/driver/postgres v1.3.10 // indirect
	gorm.io/gorm v1.23.10 // indirect
)

replace jojonomic/utils => ../../utils
