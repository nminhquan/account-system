# Use goreman to run `go get github.com/mattn/goreman`
raftgroup1-1: go run main.go --mode rm --clusterid 1 --id 1 --cluster http://172.24.20.67:12379,http://172.24.20.67:22379,http://172.24.20.67:32379 --raftport :12379 --rmport :50041 --db mas1
raftgroup1-2: go run main.go --mode rm --clusterid 1 --id 2 --cluster http://172.24.20.67:12379,http://172.24.20.67:22379,http://172.24.20.67:32379 --raftport :22379 --rmport :50042 --db mas2
raftgroup1-3: go run main.go --mode rm --clusterid 1 --id 3 --cluster http://172.24.20.67:12379,http://172.24.20.67:22379,http://172.24.20.67:32379 --raftport :32379 --rmport :50043 --db mas3
