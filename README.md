# Calendar Proxy

Proxy that changes calendar events in a pre-defined way

# TODO

## Frontend

- Use HTMX so that you can issue a PATCH

## DB

- Validate calendar before inserting :)
- Only open db once, store handle somewhere
- Write tests

## Server

- Set trusted proxies https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetHTMLTemplate
- Return html pages when input is malformed
- Return webcal link and ics file in an html page
- Write tests

## Calendar

- Write tests
