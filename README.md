# Henki

This service supports authentication through password and external login through google.

DB is in Postgres. It also requires Redis

## Unit Testing

```
./run.sh test local
```
to run all go tests

## Integration testing
first start the service:
```
./run.sh run dev
```
Then run the JS based api testing
```
cd testing/api
yarn test-local
```