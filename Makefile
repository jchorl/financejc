default:
	docker build -t jchorl/financejc .

serve:
	docker run -it --rm -p 8080:8080 -p 8000:8000 jchorl/financejc