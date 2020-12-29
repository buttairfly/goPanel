# goPanel

golang pixel panel software

## dependencies

install depdendencies

```sh
go get ./...
```

update dependencies

```sh
go get -u ./...
```

## deployment

deploy on a ledpix-raspberry-pi

```sh
./scripts/deploy.sh
```

run locally

```sh
./scripts/run-local.sh
```

## config

make sure, the config is generated/edited properly

see unit tests for configs

## unit tests

run unit tests

```sh
go test ./...
```

generate new test compare files on an error

```sh
TEST_RECORD=1 go test ./...
```

## ledpanel webfrontend

ledpanel in subfolder `gopanel/ledpanel` should get cloned from github

```sh
cd ledpanel
git clone github.com/buttairfly/ledpanel
cd ..
```

see ledpanel README.md for more details
