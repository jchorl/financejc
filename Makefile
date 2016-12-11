all: certs network db es build ui build-nginx serve nginx
dev: network db es build serve-dev ui build-nginx nginx-dev

POSTGRES_USER ?= postgres

deploy:
	$(MAKE) build
	$(MAKE) ui
	$(MAKE) build-nginx
	-docker rm -f financejc
	$(MAKE) serve
	-docker rm -f financejcnginx
	$(MAKE) nginx

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
		--restart=always \
		--network financejcnet \
		-h financejcdb \
		--expose=5432 \
		-v $(PWD)/db:/docker-entrypoint-initdb.d \
		-v financejcpgdata:/var/lib/postgresql/data \
		-e POSTGRES_USER \
		-e POSTGRES_PASSWORD \
		postgres

es: network
	docker ps | grep financejces || \
		docker run -d \
		--name financejces \
		--restart=always \
		--network financejcnet \
		-h financejces \
		--expose=9200 \
		-e ES_JAVA_OPTS="-Xms500m -Xmx500m" \
		-v financejcesdata:/usr/share/elasticsearch/data \
		elasticsearch

nginx: network
	docker ps | grep financejcnginx || \
		docker run -d \
		--name financejcnginx \
		--restart=always \
		--network financejcnet \
		-e DOMAIN=finance.joshchorlton.com \
		-v financejcletsencrypt:/etc/letsencrypt \
		-v wellknown:/usr/share/nginx/wellknown \
		-p 80:80 \
		-p 443:443 \
		jchorl/financejcnginx

nginx-dev: network
	docker ps | grep financejcnginx || \
		docker run -d \
		--name financejcnginx \
		--restart=always \
		--network financejcnet \
		-e DEV=1 \
		-e DOMAIN=finance.joshchorlton.com \
		-v $(PWD)/client/dest:/usr/share/nginx/html:ro \
		-v financejcletsencrypt:/etc/letsencrypt \
		-v wellknown:/usr/share/nginx/wellknown \
		-p 80:80 \
		-p 443:443 \
		jchorl/financejcnginx

serve: network
	docker run -d \
		--name financejc \
		--restart=always \
		--network financejcnet \
		--expose=443 \
		-h financejc \
		-e DOMAIN=finance.joshchorlton.com \
		-e PORT=443 \
		-e JWT_SIGNING_KEY \
		-e DB_ADDRESS \
		jchorl/financejc

serve-dev: network
	docker run -d \
		--name financejc \
		--restart=always \
		--network financejcnet \
		--expose=443 \
		-h financejc \
		-e DOMAIN=localhost \
		-e PORT=443 \
		jchorl/financejc

build-nginx:
	docker build -t jchorl/financejcnginx -f nginx/Dockerfile .

build:
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
	-docker volume rm financejcpgdata
	-docker volume rm financejcesdata
	-docker network rm wellknown
	-docker network rm financejcnet
	-rm client/dest/bundle.js

certs:
	docker run -it --rm \
		--name certbot \
		-v financejcletsencrypt:/etc/letsencrypt \
		-v financejcletsencryptvar:/var/lib/letsencrypt \
		-p 443:443 \
		-p 80:80 \
		quay.io/letsencrypt/letsencrypt:latest \
		certonly --standalone --noninteractive --agree-tos --keep --expand -d finance.joshchorlton.com --email=josh@joshchorlton.com


# useful targets for dev
# npm target makes it easy to add new npm packages
npm:
	docker run -it --rm \
		-v $(PWD)/client:/usr/src/app \
		-w /usr/src/app \
		node:latest /bin/bash

# connect-db connects to the postgres instance
connect-db:
	docker run -it --rm \
		--network financejcnet \
		postgres \
		psql -h financejcdb -U '$(POSTGRES_USER)'

# golang makes it easy to use tools like godep
golang:
	docker run -it --rm \
		-v $(PWD):/go/src/github.com/jchorl/financejc \
		-w /go/src/github.com/jchorl/financejc \
		golang \
		bash

# kibana makes it easy to view es data
kibana:
	docker run -it --rm \
		--name kibana \
		--network financejcnet \
		-e ELASTICSEARCH_URL=http://financejces:9200 \
		-p 5601:5601 \
		kibana

.PHONY: all ui ui-watch network dev db es nginx serve serve-dev build build-nginx clean npm connect-db golang
