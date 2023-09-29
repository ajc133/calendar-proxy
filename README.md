# Calendar Proxy

Proxy that changes calendar events in a pre-defined way

# TODO

- rename repo to calendarProxy

## DB

- Create table if not exists at server start
- Write tests

## Server

- Why do handlers have to return after calling e.g. 'c.JSON'
- Move to basic net/http server
- Set trusted proxies https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetHTMLTemplate
- Better error message when input is malformed
- Write tests

## App

- Fetch calendar
- Parse calendar
- Write tests
