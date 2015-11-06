IMAGE = fgrehm/alpine-go-web:1.5.1

.PHONY: hack
hack:
	docker run -ti --rm -v `pwd`:/code $(IMAGE)
