# Use goreman to run `go get github.com/mattn/goreman`
raftgroup1-1: go run main.go --clusterid 1 --id 1 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379 --port :50041 --db mas1
raftgroup1-2: go run main.go --clusterid 1 --id 2 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379 --port :50042 --db mas2
raftgroup1-3: go run main.go --clusterid 1 --id 3 --cluster http://127.0.0.1:12379,http://127.0.0.1:22379,http://127.0.0.1:32379 --port :50043 --db mas3

raftgroup2-1: go run main.go --clusterid 2 --id 1 --cluster http://127.0.0.1:12380,http://127.0.0.1:22380,http://127.0.0.1:32380 --port :50051 --db mas1
raftgroup2-2: go run main.go --clusterid 2 --id 2 --cluster http://127.0.0.1:12380,http://127.0.0.1:22380,http://127.0.0.1:32380 --port :50052 --db mas2
raftgroup2-3: go run main.go --clusterid 2 --id 3 --cluster http://127.0.0.1:12380,http://127.0.0.1:22380,http://127.0.0.1:32380 --port :50053 --db mas3

raftgroup3-1: go run main.go --clusterid 3 --id 1 --cluster http://127.0.0.1:12381,http://127.0.0.1:22381,http://127.0.0.1:32381 --port :50061 --db mas1
raftgroup3-2: go run main.go --clusterid 3 --id 2 --cluster http://127.0.0.1:12381,http://127.0.0.1:22381,http://127.0.0.1:32381 --port :50062 --db mas2
raftgroup3-3: go run main.go --clusterid 3 --id 3 --cluster http://127.0.0.1:12381,http://127.0.0.1:22381,http://127.0.0.1:32381 --port :50063 --db mas3



