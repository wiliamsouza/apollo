Web client
==========

Test using [React](http://facebook.github.io/react/index.html)
to write the Apollo web api client.

Usage
-----

It uses [npm](https://npmjs.org/) and [Bower](http://bower.io/)
to handler dependencies.

Install [Node.js](http://nodejs.org/download/) before run the commands below.

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
go build ../bin/apollod.go
```
Start RESTful server.

``` 
../bin/apollod -config ../etc/apollod.conf &
```

Start web api test server:

```
node start
```

Point your browser to http://localhost:3000/.

Development
-----------

When making modifcation to React `JSX` file run first:

```
jsx --watch scripts/src/ scripts/build/
```
