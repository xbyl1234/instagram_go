module makemoney

go 1.16

replace golang.org/x/net => ./../third/net

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.8.0
	github.com/edgeware/mp4ff v0.27.0
	github.com/emersion/go-imap v1.2.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/klauspost/compress v1.15.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/twmb/murmur3 v1.1.6
	github.com/utahta/go-cronowriter v1.2.0
	go.mongodb.org/mongo-driver v1.9.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad
)
