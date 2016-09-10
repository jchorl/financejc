all: ui build serve

ui:
	docker run --rm --name uibuild -it -v $(PWD):/usr/src/app -w /usr/src/app/client node:latest /bin/bash -c "npm install; node ./node_modules/.bin/webpack --progress --colors"
ui-watch:
	docker run --rm --name uibuild -it -v $(PWD):/usr/src/app -w /usr/src/app/client node:latest /bin/bash -c "npm install; node ./node_modules/.bin/webpack --progress --colors --watch"
serve:
	docker run -it --rm -p 8080:8080 -e PORT 8080 jchorl/financejc
build:
	docker build -t jchorl/financejc .
.PHONY: all ui ui-watch serve build
