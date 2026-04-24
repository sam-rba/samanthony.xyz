MODS = $(shell git submodule | awk '{print $$2}')

build: fmt ${MODS} static/resume.png
	hugo build
	find public -name '*.html' -exec tidy -mqi {} \;

fmt:
	gotmplfmt -w .

publish: build
	rsync -rz --delete public/ samanthony.xyz:/var/www/htdocs/samanthony.xyz/
	$(eval CID := $(shell ipfs add -r public | tail -n1 | awk '{print $$2}')) \
	ssh m900 "IPFS_PATH=/var/kubo ipfs pin add -r ${CID}"
	ipfs name publish --key=www.samanthony.xyz ${CID}

serve: build
	hugo serve -DEF

clean:
	rm -rf public

static/resume.png: static/resume.pdf
	gs -dNOPAUSE -dBATCH -sDEVICE=png16m -r300 -sOutputFile=$@ $<

${MODS}:
	git submodule update --init --recursive

.PHONY: build fmt publish serve clean
