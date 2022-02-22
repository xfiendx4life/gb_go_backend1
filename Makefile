docker_dir := ./storage/docker_runner
compose_dir := ./
init_db: $(docker_dir)
	sudo bash $(docker_dir)/restart.sh 
	sleep 3s 
	sudo bash $(docker_dir)/migrate_up.sh
test-unit:
	go test ./... -coverprofile=cover.out
test-integration: $(docker_dir)
	sudo bash $(docker_dir)/test/restart_test.sh 
	sleep 3s 
	sudo bash $(docker_dir)/test/migrate_test_up.sh
	go test -tags=integration ./... -coverprofile cover.out
	sudo docker stop postgres_test && sudo docker rm postgres_test
run:
	make init_db
	TIMEOUT=2 PORT=8080 LOGLEVEL=debug DATABASE_URL=postgres://xfiendx4life:123456@172.17.0.2:5432/shortener MAXCONS=10 MINCONS=5 SECRETKEY="somesecret" TTL=60 go run cmd/shrtener/main.go
run-docker-full:
	sudo rm -rf $(docker_dir)/_data
	sudo docker-compose build --no-cache
	sudo docker-compose up
run-docker:
	sudo docker-compose up
