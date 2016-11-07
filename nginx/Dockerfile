FROM nginx

ADD client/dest/* /usr/share/nginx/html/
ADD scripts/renew_certs.sh /etc/cron.daily/renew_certs.sh
ADD nginx/nginx.conf /etc/nginx/nginx.conf
ADD nginx/nginx_start.sh /financejc/

CMD /financejc/nginx_start.sh
