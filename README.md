# Calendar Proxy

Proxy that changes calendar events in a pre-defined way

# TODO

## DB

- Update a row
- Refresh calendarBody daily for all records
- Write tests

## Server

- Set trusted proxies https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetHTMLTemplate
- Return html pages when input is malformed
- PATCH a calendar's replacementSummary
    - Clear cache (maybe only for this key?) upon doing so
- Pass a GET query option that bypasses cache
- Write tests

## Calendar

- Write tests
