FROM nginx

# install openssl for cert generation for dev env :(
RUN apt-get update && apt-get install -y openssl

ADD scripts/renew_certs.sh /etc/cron.daily/renew_certs.sh
ADD nginx/nginx.conf /etc/nginx/nginx.conf
ADD nginx/nginx_start.sh /financejc/
ADD client/build /usr/share/nginx/html/

CMD /financejc/nginx_start.sh
