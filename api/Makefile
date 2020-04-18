pkg = github.com/marsDev31/kuclap-backend/api
APP = kuclap-api
DYNO = web
uppkg: push-and-get-go
goget: go-get
herokucp: heroku-container-push-web
herokucr: heroku-container-release-web
hpr: cp-cr

go-get:
	go get -u ${pkg}

push-and-get-go:
	git add . && git commit -m '[+] update : pkg on version control' && git push && go get -u ${pkg}

heroku-login:
	heroku login

heroku-container-login:
	heroku container:login

heroku-container-push-web:
	heroku container:push ${DYNO} --app ${APP}

heroku-container-release-web:
	heroku container:release ${DYNO} --app ${APP}

cp-cr:
	heroku container:push ${DYNO} --app ${APP} && heroku container:release ${DYNO} --app ${APP}