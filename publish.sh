prepare() {
	git config --global user.email "libsora25@gmail.com"
	git config --global user.name "travis-ci"
	git config push.default simple
	mkdir -p .travis
	openssl aes-256-cbc -K $encrypted_d7a476c7b0de_key -iv $encrypted_d7a476c7b0de_iv -in id_rsa.enc -out .travis/id_rsa -d
	chmod 600 .travis/id_rsa
	eval `ssh-agent -s`
	ssh-add .travis/id_rsa
	virtualenv .venv
	. ./.venv/bin/activate
	pip install ghp-import
}

publish() {
	git remote remove origin
	git remote add origin git@github.com:if1live/umi.git
	git fetch origin
	ghp-import output
	git push origin gh-pages
}
