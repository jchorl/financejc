ui:
	docker run --rm --name uibuild -it -v $(PWD):/usr/src/app -w /usr/src/app/client node:latest /bin/bash -c "npm install; node ./node_modules/.bin/webpack --progress --colors"
ui-watch:
	docker run --rm --name uibuild -it -v $(PWD):/usr/src/app -w /usr/src/app/client node:latest /bin/bash -c "npm install; node ./node_modules/.bin/webpack --progress --colors --watch"