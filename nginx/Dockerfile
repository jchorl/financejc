FROM node
WORKDIR /usr/src/app
ADD client/ .
RUN yarn
RUN yarn run build
RUN gzip --best -f --keep build/static/js/*
RUN gzip --best -f --keep build/static/css/*

FROM nginx
# install openssl for cert generation for dev env :(
RUN apt-get update && apt-get install -y openssl
ADD nginx/nginx.conf /etc/nginx/nginx.conf
ADD nginx/nginx_start.sh /financejc/
COPY --from=0 /usr/src/app/build /usr/share/nginx/html/
CMD /financejc/nginx_start.sh
