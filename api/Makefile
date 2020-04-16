pkg = github.com/marsDev31/kuclap-backend/api
uppkg: push-and-get-go
goget: go-get
herokucp: heroku-container-push-web
herokucr: heroku-container-release-web:

go-get:
	go get -u ${pkg}

push-and-get-go:
	git add . && git commit -m '[+] update : pkg on version control' && git push && go get -u ${pkg}

heroku-login:
	heroku login

heroku-container-login:
	heroku container:login

heroku-container-push-web:
	heroku container:push web --app kuclap-api

heroku-container-release-web:
	heroku container:release web --app kuclap-api