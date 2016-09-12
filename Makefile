all: network db ui build serve

network:
	docker network ls | grep financejcnet || docker network create financejcnet
ui: client/dest/bundle.js;
client/dest/bundle.js: $(shell find client/src)
	docker run --rm --name uibuild -it -v $(PWD)/client:/usr/src/app -w /usr/src/app node:latest /bin/bash -c "npm install; node ./node_modules/.bin/webpack --progress --colors"
ui-watch:
	docker run --rm --name uiwatch -it -v $(PWD)/client:/usr/src/app -w /usr/src/app node:latest /bin/bash -c "npm install; node ./node_modules/.bin/webpack --progress --colors --watch"
db: network
	docker ps | grep financejcdb || docker run --name financejcdb --network financejcnet -h financejcdb --expose=5432 -v $(PWD)/db:/docker-entrypoint-initdb.d -e POSTGRES_USER=financejc -e POSTGRES_PASSWORD=financejc -d postgres
connect-db:
	docker run -it --rm --network financejcnet postgres psql -h financejcdb -U financejc
serve: network
	docker run -it --name financejc --rm --network financejcnet -p 4443:443 -e PORT=443 -v $(PWD)/client/dest:/go/src/github.com/jchorl/financejc/client/dest jchorl/financejc
build: ui
	docker build -t jchorl/financejc .
.PHONY: all ui ui-watch network db connect-db serve build
