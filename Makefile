
gitadd:
	git remote add github https://github.com/fmnx-su/pack
	git remote add codeberg https://codeberg.org/fmnx/pack

push:
	git push
	git push github
	git push codeberg

build:
	pack -Bqs