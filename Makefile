all: network db es build serve nginx
dev: network db es build serve-dev nginx

network:
	docker network ls | grep financejcnet || \
		docker network create financejcnet

ui: client/dest/bundle.js;
client/dest/bundle.js: $(shell find client/src) client/webpack.production.config.js
	docker run -it --rm \
		--name uibuild \
		-v $(PWD)/client:/usr/src/app \
		-w /usr/src/app \
		node:latest \
		/bin/bash -c "npm install; NODE_ENV=production node ./node_modules/.bin/webpack -p --config webpack.production.config.js --progress --colors"

ui-watch:
	docker run -it --rm \
		--name uiwatch \
		-v $(PWD)/client:/usr/src/app \
		-w /usr/src/app \
		node:latest \
		/bin/bash -c "npm install; node ./node_modules/.bin/webpack --progress --colors --watch"

db: network
	docker ps | grep financejcdb || \
		docker run -d \
		--name financejcdb \
		--network financejcnet \
		-h financejcdb \
		--expose=5432 \
		-v $(PWD)/db:/docker-entrypoint-initdb.d \
		-v pgdata:/var/lib/postgresql/data \
		-e POSTGRES_USER=financejc \
		-e POSTGRES_PASSWORD=financejc \
		postgres

es: network
	docker ps | grep financejces || \
		docker run -d \
		--name financejces \
		--network financejcnet \
		-h financejces \
		--expose=9200 \
		-v esdata:/usr/share/elasticsearch/data \
		elasticsearch

nginx: network build-nginx
	docker ps | grep financejcnginx || \
		docker run -d \
		--name financejcnginx \
		--network financejcnet \
		-e DEV=1 \
		-e DOMAIN=finance.joshchorlton.com \
		-v $(PWD)/client/dest:/usr/share/nginx/html:ro \
		-v letsencrypt:/etc/letsencrypt \
		-p 4443:443 \
		jchorl/financejcnginx

serve: network
	docker run -d \
		--name financejc \
		--network financejcnet \
		--expose=4443 \
		-h financejc \
		-e DOMAIN=finance.joshchorlton.com \
		-e PORT=443 \
		-v $(PWD)/client/dest:/go/src/github.com/jchorl/financejc/client/dest \
		jchorl/financejc

serve-dev: network
	docker run -d \
		--name financejc \
		--network financejcnet \
		--expose=4443 \
		-h financejc \
		-e DOMAIN=localhost \
		-e PORT=4443 \
		-v $(PWD)/client/dest:/go/src/github.com/jchorl/financejc/client/dest \
		jchorl/financejc

build-nginx:
	docker build -t jchorl/financejcnginx -f nginx/Dockerfile .

build: ui
	docker build -t jchorl/financejc .

test: clean network db es test-all
test-all:
	docker run -it --rm \
		--name financejctest \
		--network financejcnet \
		-v $(PWD):/go/src/github.com/jchorl/financejc \
		-w /go/src/github.com/jchorl/financejc \
		golang \
		sh -c 'go test --tags=integration $$(go list ./... | grep -v /vendor/)'

clean:
	-docker rm -f financejcdbcon
	-docker rm -f financejcdb
	-docker rm -f financejces
	-docker rm -f financejcnginx
	-docker rm -f financejc
	-docker network rm financejcnet
	-rm client/dest/bundle.js

npm:
	docker run -it --rm \
		-v $(PWD)/client:/usr/src/app \
		-w /usr/src/app \
		node:latest /bin/bash

connect-db:
	docker run -it --rm \
		--network financejcnet \
		postgres \
		psql -h financejcdb -U financejc

golang:
	docker run -it --rm \
		-v $(PWD):/go/src/github.com/jchorl/financejc \
		-w /go/src/github.com/jchorl/financejc \
		golang \
		bash

kibana:
	docker run -it --rm \
		--name kibana \
		--network financejcnet \
		-e ELASTICSEARCH_URL=http://financejces:9200 \
		-p 5601:5601 \
		kibana

.PHONY: all ui ui-watch network dev db es nginx serve serve-dev build build-nginx clean npm connect-db golang
