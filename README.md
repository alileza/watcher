# Watcher

Simple binary that will watch your log files & notify for any words that match specification.

## Usage
### Instalation
```sh
go get github.com/alileza/watcher
```
Make sure that you add `$GOPATH/bin` as one of your path
### Run apps
```sh
watcher -path '/var/log/test.log' \
 -words 'hello,world,hehe,mantap' \
 -webhook 'https://hooks.slack.com/services/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' \
 -channel 'tools-err'
```
