pkg = github.com/KUclap/KUclap-review
# APP_PROD = kuclap-api
APP_STAGING = kuclap-api-staging
DYNO = web
uppkg: push-and-get-go
goget: go-get
herokucp: heroku-container-push-web
herokucr: heroku-container-release-web
# deploy-to-prod: cp-cr-prod
deploy-to-staging: cp-cr-staging

### Development (PRODUCTION)

gomodd:
	modd -f ./config/modd.conf 
	
run-on-prod:
	KIND=production ./builder/kuclap-review-api
	
go-get:
	go get -u ${pkg}

push-and-get-go:
	git add . && git commit -m '[+] update : pkg on version control' && git push && go get -u ${pkg}

heroku-login:
	heroku login

heroku-container-login:
	heroku container:login

# heroku-container-push-web:
# 	heroku container:push ${DYNO} --app ${APP_PROD}

# heroku-container-release-web:
# 	heroku container:release ${DYNO} --app ${APP_PROD}

### Deploy to Heroku
# cp-cr-prod:
# 	heroku container:push ${DYNO} --app ${APP_PROD} && heroku container:release ${DYNO} --app ${APP_PROD}

cp-cr-staging:
	heroku container:push ${DYNO} --app ${APP_STAGING} && heroku container:release ${DYNO} --app ${APP_STAGING}
# heroku container:push ${DYNO} --recursive --context-path . --app ${APP_STAGING} && heroku container:release ${DYNO} --app ${APP_STAGING}
# NOTE: Dockerfile's name should starting with 'D' capital letter, then use --recursive optional for recursive search the image file and use context-path optional for building context on root.

### Load Testing 
load-testing-vegeta:
	vegeta attack -duration=1s -rate=1000 -targets=./config/vegeta.conf | vegeta report

### MongoDB 
duplicate-collection-to-another-db:
	### shell to mongo cluster
	use src_db
	db.src_collection.find().forEach(function(d){ db.getSiblingDB('dest_db')['dest_collection'].insert(d); });

### Build docker container
build-docker:
	docker build -f ./docker/Prod.Dockerfile -t kuclap-review-api-v2
