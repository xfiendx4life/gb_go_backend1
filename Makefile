docker_dir := ./storage/docker_runner
init_db: $(docker_dir)
	sudo bash $(docker_dir)/restart.sh 
	sleep 3s 
	sudo bash $(docker_dir)/migrate_up.sh
test-unit:
	go test -tags=unit ./... -coverprofile=cover.out
test-integration: $(docker_dir)
	sudo bash $(docker_dir)/test/restart_test.sh 
	sleep 3s 
	sudo bash $(docker_dir)/test/migrate_test_up.sh
	go test -tags=integration ./... -coverprofile cover.out
	sudo docker stop postgres_test && sudo docker rm postgres_test
run:
	make init_dbko
	go run cmd/shrtener/main.go
