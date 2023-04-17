run:
	docker rm -f searcher
	docker build . -t searcher
	docker run --name searcher -p 8000:8000 -p 8001:8001 searcher