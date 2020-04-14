pkg = github.com/marsDev31/kuclap-backend/api
uppkg: push-and-get-go
goget: go-get

go-get:
	go get -u ${pkg}

push-and-get-go:
	git add . && git commit -m '[+] update : pkg on version control' && git push && go get -u ${pkg}
