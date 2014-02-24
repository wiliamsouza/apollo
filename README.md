Apollo
======

Easy way to write, build, run and collect test results.

Install
-------

Install Go [tools](http://golang.org/doc/install#install) I prefer use
[godeb](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go) to do that.

Get the source code:

```
git clone https://github.com/wiliamsouza/apollo.git
```

Go to `bin` source directory and build the source:

```
cd apollo/bin
go build apollod.go
```

Start the webserver:
 
The only configuration that need to be adjusted is keys path
the following example set the path to a test keys.

```
rsa:
  private: "../data/keys/rsa"
  public: "../data/keys/rsa.pub"
```

```
./apollod -config ../etc/apollod.conf
```

Now you can test the API.

User API
--------

Adding a new `user`:

```
curl -v --header "Content-Type: application/json" --request POST --data '{"name":"Jhon Doe","email":"jhon@doe.com","password":"12345"}' http://localhost:8000/users
```

It will return `user` name and email:

```
{"name":"Jhon Doe","email":"jhon@doe.com"}
```

Authenticate `user`:

```
curl -v --header "Content-Type: application/json" --request POST --data '{"email":"jhon@doe.com","password":"12345"}' http://localhost:8000/users/authenticate
```
It will return a JWT(JSON Web Token) `token`:

```
{"token":"cCI6IkpXVCJCI6Impob25AZG9lLmNvbSIsImV4cCI6MTM5MzMxODU0OH0ZLoU"}
```

From now you should add and `Authorization` header like:

```
--header "Authorization: Bearer <token>"
```

`<token>` is the one returned by `users/authenticate`.

Detail `user`:

```
curl -v --header "Content-Type: application/json" --header "Authorization: Bearer <token>" --request GET http://localhost:8000/users/jhon@doe.com
```
It will return `user` details:

```
{"name":"Jhon Doe","email":"jhon@doe.com","apikey":"0Yy00ZDRhLThmNmQt...","created":"2014-02-22T00:20:44.511-03:00","lastlogin":"0001-01-01T00:00:00Z"}
```

Package API
-----------

Uploading a `package`:

```
curl -vv --header "Authorization: Bearer <token>"  --request POST --form package=@package1.tgz --form metadata=@metadata1.json http://localhost:8000/tests/packages
```

The `data` directory contains some files to test.

It will return basic information about `package`:

```
{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}}
```

Listing `packages`:

```
curl -vv --header "Authorization: Bearer <token>"  --header "Content-Type: application/json" -X GET http://localhost:8000/tests/packages
```

It will return a list of `packages`:

```
[{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}}]
```

Detailing `package`:

```
curl -vv --header "Authorization: Bearer <token>"  --header "Content-Type: application/json" -X GET http://localhost:8000/tests/packages/package1.tgz
```

It will return `package` detail:

```
{"filename":"package1.tgz","metadata":{"version":0.1,"description":"Package1 ON/OFF test","install":"adb push dist/package1.jar /data/local/tmp/","run":"adb shell uiautomator runtest package1.jar -c com.github.wiliamsouza.package1.Package1Test"}}
```

Downloading `package`:

```
curl -vv --header "Authorization: Bearer <token>"  -X GET http://localhost:8000/tests/packages/downloads/package1.tgz -o package.tgz
```

It will download `package`.

Organization API
----------------

Adding a new `organization`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer <token>"  --request POST --data '{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}' http://localhost:8000/organizations
```

It will return:

```
{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}
```

List `organizations`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer <token>"  --request GET http://localhost:8000/organizations
```

It will return:

```
[{"name":"janecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jane@doe.com"]},{"name":"jhoncorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}]
```

Detail `organization`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer <token>"  --request GET http://localhost:8000/organizations/doecorp
```

It will return:

```
{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}
```

Update `organization`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer <token>"  --request PUT --data '{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com"]}],"admins":["jane@doe.com"]}' http://localhost:8000/organizations/doecorp
```

Update is used to simulate teams and admins deletion, change and addtion.

Delete:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer <token>"  --request DELETE http://localhost:8000/organizations/doecorp
```

[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/wiliamsouza/apollo/trend.png)](https://bitdeli.com/free "Bitdeli Badge")
