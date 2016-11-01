all: network db ui build serve
dev: network db ui build serve-dev

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
		-e POSTGRES_USER=financejc \
		-e POSTGRES_PASSWORD=financejc \
		postgres

serve: network
	docker run -it --rm \
		--name financejc \
		--network financejcnet \
		-p 4443:443 \
		-e DOMAIN=finance.joshchorlton.com \
		-e PORT=443 \
		-v $(PWD)/client/dest:/go/src/github.com/jchorl/financejc/client/dest \
		jchorl/financejc

serve-dev: network
	docker run -it --rm \
		--name financejc \
		--network financejcnet \
		-p 4443:4443 \
		-e DEV=1 \
		-e DOMAIN=localhost \
		-e PORT=4443 \
		-v $(PWD)/client/dest:/go/src/github.com/jchorl/financejc/client/dest \
		jchorl/financejc

build: ui
	docker build -t jchorl/financejc .

test: clean network db test-all
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
	-docker rm -f financejc
	-docker network rm financejcnet
	-rm client/dest/bundle.js

npm:
	docker run -it --rm \
		--name npm \
		-v $(PWD)/client:/usr/src/app \
		-w /usr/src/app \
		node:latest /bin/bash

connect-db:
	docker run -it --rm \
		--name financejcdbcon \
		--network financejcnet \
		postgres \
		psql -h financejcdb -U financejc

.PHONY: all ui ui-watch network db serve build clean npm connect-db
