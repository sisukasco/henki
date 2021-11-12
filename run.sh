#!/bin/bash
# run.sh db reset
# run.sh run local
# run.sh migrate up
# run.sh migrate down

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
    echo "Running auth in development mode ..."
    SISUKAS_ENVIRONMENT=development \
    SISUKAS_LOGGED_IN_USER_FOR_TESTING=testing_user \
    go run ./cmd --conf conf-dev.yaml serve
}

function runStaging(){
    echo "Running auth in staging mode ..."
    
    go run ./cmd --conf conf-staging.yaml serve
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
    test)
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

