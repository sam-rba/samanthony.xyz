MODS = $(shell git submodule | awk '{print $$2}')

build: fmt ${MODS}
	hugo build
	find public -name '*.html' -exec tidy -mqi {} \;

fmt:
	gotmplfmt -w .

publish: build
	rsync -rz --delete public/ samanthony.xyz:/var/www/htdocs/samanthony.xyz/

serve: build
	hugo serve -DEF

clean:
	rm -rf public

${MODS}:
	git submodule update --init --recursive

.PHONY: build fmt publish serve clean
