Request Catcher
===============

[Request Catcher](http://requestcatcher.com) is a tool for catching web requests for testing webhooks, http clients and other applications that communicate over http. Request Catcher gives you a subdomain to test your application against. Keep the index page open and instantly see all incoming requests to the subdomain via WebSockets.

### Persistence

Request Catcher does not currently persist requests. You will only receive requests that are sent while you have the listening page open. Requests before or after will be lost.

### Running locally

To run Request Catcher locally, ensure that you've installed Go and all the project dependencies. You'll need to create the MySQL database and apply migrations:

```
mysql -e "CREATE DATABASE requestcatcher_development;"
goose up
```

When starting the server, the command line interface takes one argument, the path to a json configuration file. See `config/development.json` to see the possible configuration parameters.

`go run main.go "config/development.json"`

Then visit `http://lvh.me:8080` in your browser.

