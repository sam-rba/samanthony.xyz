build:
	hugo build

publish: build
	rsync -rz --delete public/ samanthony.xyz:/var/www/htdocs/samanthony.xyz/

serve:
	hugo serve -DEF

clean:
	rm -rf public
