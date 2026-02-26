build: static/resume.png
	hugo build

publish: build
	rsync -rz --delete public/ samanthony.xyz:/var/www/htdocs/samanthony.xyz/
	ipfs add -r public \
		| tail -n1 | awk '{print $$2}' \
		| xargs -I{} ssh samanthony.xyz "IPFS_PATH=/var/kubo ipfs pin add -r {}"

serve:
	hugo serve -DEF

clean:
	rm -rf public

static/resume.png: static/resume.pdf
	gs -dNOPAUSE -dBATCH -sDEVICE=png16m -r300 -sOutputFile=$@ $<
