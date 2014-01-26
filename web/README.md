Web client
==========

Test using [react](http://facebook.github.io/react/index.html)
to write the Apollo web api client.

Usage
-----

It uses [npm](https://npmjs.org/) and [bower](http://bower.io/)
to handler dependencies.

Install [node.js](http://nodejs.org/download/) before run the commands below.

```
npm install
```

The above command will download and install node.js dependencies.

```
bower install
```

The above command will download and install project dependencies.

First you need build `bin/apollo-webserver.go` before start test server:

```
go build ../bin/apollo-webserver.go
```
Start RESTful server.

``` 
../bin/apollo-webserver -config ../etc/apollo.conf &
```

Start web api test server:

```
node start
```

Point your browser to http://localhost:3000/.
