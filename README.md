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

When starting the server, the command line interface takes three arguments: the hostname to listen on, the port number to listen on and the hostname to consider as the 'root'. The root hostname will serve the front page with information about the application. Only subdomains of the root hostname or other hosts routed to the application will catch requests. When running locally, you can use the `lvh.me` domain to test subdomains. For example, launch as

`go run main.go start localhost 8080 lvh.me`

Then visit `http://lvh.me:8080` in your browser.

