docker run -it --rm --name certbot \
	-v financejcletsencrypt:/etc/letsencrypt \
	-v wellknown:/wellknown \
	quay.io/letsencrypt/letsencrypt:latest \
	certonly -a webroot --webroot-path=/wellknown --noninteractive --agree-tos --keep --expand -d finance.joshchorlton.com --email=josh@joshchorlton.com
