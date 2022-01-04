docker_dir := ./storage/docker_runner
init_db: $(docker_dir)
	sudo bash $(docker_dir)/restart.sh 
	sleep 2s 
	sudo bash $(docker_dir)/migrate_up.sh
