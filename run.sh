#!/bin/bash
# run.sh db reset
# run.sh run local
# run.sh migrate up
# run.sh migrate down
set -o allexport
if [ -f .env ]
then
    source .env
fi
set +o allexport

function resetDevDB(){
    echo "reseting the Dev DB..."
    resetDB "authdb_dev"
    echo "migrating ..."
    migrateDevDBUp
}

function resetStagingDB(){
    echo "reseting the Staging DB..."
    resetDB "authdb_staging"
    echo "migrating ..."
    migrateStagingDBUp
}

function resetDB(){
    
    HENKI_DB_NAME=$1
    
    psql -h localhost -U dockuser postgres << SQL
 drop database IF EXISTS $HENKI_DB_NAME;
 
 CREATE DATABASE "$HENKI_DB_NAME"
    ENCODING 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TEMPLATE template0;
    
 grant all privileges on database $HENKI_DB_NAME to dockuser;
SQL
}

function migrateDevDBUp(){
    go run ./cmd --source ./pkg/db/migrations  --conf conf-dev.yaml migrate up 
}

function migrateDevDBDown(){
    go run ./cmd --source ./pkg/db/migrations  --conf conf-dev.yaml migrate reset 
}

function migrateStagingDBUp(){
    go run ./cmd --source ./pkg/db/migrations  --conf conf-staging.yaml migrate up 
}

function migrateStagingDBDown(){
    go run ./cmd --source ./pkg/db/migrations  --conf conf-staging.yaml migrate reset 
}

function runDev(){
    if [ ! -f "conf-dev.yaml" ]; then
        echo "Configuration file conf-dev.yaml does not exist."
        return 1
    fi    
    echo "Running auth in development mode ..."
    SISUKAS_ENVIRONMENT=development \
    SISUKAS_LOGGED_IN_USER_FOR_TESTING=testing_user \
    go run ./cmd --conf conf-dev.yaml serve
}

function runStaging(){
    if [ ! -f "conf-staging.yaml" ]; then
        echo "Configuration file conf-staging.yaml does not exist."
        return 1
    fi    
    echo "Running auth in staging mode ..."
    
    go run ./cmd --conf conf-staging.yaml serve
}

function buildLinuxVersion(){
    mkdir -p ./deploy/bin
    rm ./deploy/bin/henki
    perl -pe '/^VERSION=/ and s/(\d+\.\d+\.)(\d+)/$1 . ($2+1)/e' -i ./version
    source ./version
    #VERSION=1.0.17
    echo "Compiling Henki ..."
    BUILD_TIME=$(date)
    
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/sisukasco/henki/pkg/version.BuildTime=$BUILD_TIME' -X 'github.com/sisukasco/henki/pkg/version.Version=$VERSION'" -a -installsuffix cgo -o ./deploy/bin/henki ./cmd
    
    test $? -eq 0 || exit 1    
}

function runAllTests(){
    SISUKAS_ENVIRONMENT=development \
    SISUKAS_LOGGED_IN_USER_FOR_TESTING=testing_user \
    go test -p 1 -count=1 $(go list ./...)
}

function runTests(){
    echo "Running tests in folder $1 ..."
    SISUKAS_ENVIRONMENT=development \
    SISUKAS_LOGGED_IN_USER_FOR_TESTING=testing_user \
    go test -v -count=1 $1
}

function runTestFn(){
    echo "Running test $1 in folder $2 ..."
    SISUKAS_ENVIRONMENT=development \
    SISUKAS_LOGGED_IN_USER_FOR_TESTING=testing_user \
    go test -v -count=1 -run $1 $2
}

if (( $# < 2 ))
then
      echo "Usage cmd.sh command what"
      exit 1
fi

COMMAND=$1
WHAT=$2

case $COMMAND in
    resetdb)
        case $WHAT in
            dev)
                resetDevDB
            ;;
            staging)
                resetStagingDB
            ;;
        esac
    ;;
    migrate)
        case $WHAT in
            dev)
                migrateDevDBUp
            ;;
            staging)
                migrateStagingDBUp
            ;;
        esac
    ;;
    run)
        case $WHAT in
            dev)
                runDev
            ;;
            staging)
                runStaging
            ;;
        esac
    ;;
    build)
        case $WHAT in
            linux-version)
                buildLinuxVersion
            ;;
        esac
    ;;
    test)
        if [ ! -f "conf-dev.yaml" ]; then
            echo "Error: Configuration file conf-dev.yaml does not exist."
            exit 1
        fi
        case $WHAT in
           local)
                runAllTests
           ;;
           fn)
                runTestFn $3 $4
           ;;            
           *)
                runTests $WHAT
           ;;
        esac
        
    ;;    
esac

