# for Docker Desktop on WSL2
host_ip := $(shell ip -o route get 8.8.8.8 | perl -lne 'print $$1 if /via (\S+)/')

dev:
	tmux \
		new make up \;\
		split make -j somux unison/watch \;\
		set synchronize-panes \;\
		select-layout even-v

up:
	docker-compose up

somux: wait/app wait/unison wait/somux
	somux \
		-R 0.0.0.0:9000:${host_ip}:9000 \
		-R 0.0.0.0:9003:${host_ip}:9003 \
		-L 9876:app:9876 \
		-L 5000:unison:5000 \
		docker-compose exec -T somux somux --

unison := \
	unison . socket://localhost:5000//code \
		-ignore 'Name .git' \
		-ignore 'Name .idea' \
		-ignore 'Path vendor' \
		-ignore 'Path node_modules'

unison: wait/5000
	${unison} -auto

unison/watch: wait/5000
	${unison} -auto -batch -repeat watch

wait/5000:
	for i in {1..30}; do nc -z -v -w 1 localhost 5000 && break; sleep 0.5; done

wait/%:
	for i in {1..30}; do docker-compose exec -T ${@F} /bin/true && break; sleep 0.5; done
