FROM nginx

ADD client/dest/* /usr/share/nginx/html/
ADD nginx/nginx.conf /etc/nginx/nginx.conf
ADD nginx/nginx_start.sh /financejc/

CMD /financejc/nginx_start.sh
