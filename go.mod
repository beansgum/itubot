module github.com/C-ollins/itubot

go 1.13

require (
	github.com/DataDog/zstd v1.4.8 // indirect
	github.com/Sereal/Sereal v0.0.0-20200820125258-a016b7cda3f3 // indirect
	github.com/adshao/go-binance/v2 v2.2.1
	github.com/asdine/storm v2.1.2+incompatible
	github.com/binance-exchange/go-binance v0.0.0-20180518133450-1af034307da5 // indirect
	github.com/decred/slog v1.2.0
	github.com/go-kit/kit v0.10.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jrick/logrotate v1.0.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
)

replace github.com/adshao/go-binance/v2 => github.com/C-ollins/go-binance/v2 v2.2.2-0.20210320142445-c3818201a0c5
