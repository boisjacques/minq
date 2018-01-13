export "MINQ_LOG"=mp,mutex
rm clientlog
go build -o clnt main.go
./clnt -addr=192.168.178.60:4433
