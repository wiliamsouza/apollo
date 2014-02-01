Apollo
======

Easy way to write, build, run and collect test results.

Install
-------

Install Go [tools](http://golang.org/doc/install#install) I prefer use
[godeb](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go) to do that.

Get the source code:

     git clone https://github.com/wiliamsouza/apollo.git

Go to `bin` source directory and build the source:

    cd apollo/bin
    go build apollo-webserver.go

Start the webserver:

    ./apollo-webserver -config ../etc/apollo-webserver.conf

Now you can test the API.

Package API
-----------

Uploading a `package`:

    curl -vv --request POST --form package=@package1.tgz --form metadata=@metadata1.json http://localhost:8000/tests/packages

The `data` directory contains some files to test.

It will return basic information about `package`:

    {"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}}

Listing `packages`:


    curl -vv --header "Content-Type: application/json" -X GET http://localhost:8000/tests/packages

It will return a list of `packages`:

    [{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}}]

Detailing `package`:

    curl -vv --header "Content-Type: application/json" -X GET http://localhost:8000/tests/packages/package1.tgz

It will return `package` detail:

    {"filename":"package1.tgz","metadata":{"version":0.1,"description":"Package1 ON/OFF test","install":"adb push dist/package1.jar /data/local/tmp/","run":"adb shell uiautomator runtest package1.jar -c com.github.wiliamsouza.package1.Package1Test"}}

Downloading `package`:

    curl -vv -X GET http://localhost:8000/tests/packages/downloads/package1.tgz -o package.tgz

It will download `package`.

User API
--------

Adding a new `user`:

    curl -v --header "Content-Type: application/json" --request POST --data '{"name":"Jhon Doe","email":"jhon@doe.com","password":"12345"}' http://localhost:8000/users

It will return `user` name and email:

    {"name":"Jhon Doe","email":"jhon@doe.com"}


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/wiliamsouza/apollo/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

