build: static/resume.png
	hugo build

publish: build
	rsync -rz --delete public/ samanthony.xyz:/var/www/htdocs/samanthony.xyz/

serve:
	hugo serve -DEF

clean:
	rm -rf public

static/resume.png: static/resume.pdf
	gs -dNOPAUSE -dBATCH -sDEVICE=png16m -r300 -sOutputFile=$@ $<
