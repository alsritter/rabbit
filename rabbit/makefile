.PHONY: run-server
run-server: run-export-env-variable
	go run .

run-export-env-variable:
	export ETCD_SERVER_URL=172.16.238.101:2379
	export ETCD_USERNAME=
	export ETCD_PASSWORD=