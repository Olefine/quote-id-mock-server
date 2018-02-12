version ?= `cat .version`

build:
	docker build -t "olefine/quote-id-mock:${version}" .
run:
	docker run --rm -d -p 3030:8080 "olefine/quote-id-mock:${version}"
