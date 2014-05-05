Apollo
======

[![wercker status](https://app.wercker.com/status/f3fc24f278c2ca43fce270d5534745cc/m/master "wercker status")](https://app.wercker.com/project/bykey/f3fc24f278c2ca43fce270d5534745cc)

Easy way to write, build, run and collect test results.

Install
-------

Install Go [tools](http://golang.org/doc/install#install) I prefer use
[godeb](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go) to do that.

Get the source code:

```
git clone https://github.com/wiliamsouza/apollo.git
```

After that install godeps:

```
go get github.com/tools/godep
```

Install dependencies:

```
cd apollo/Godeps
godeps restore
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

MongoDB container
-----------------

Apollo uses MongoDB as database and for development we run it inside a Docker
container to get it running follow the instructions bellow:

```
sudo apt-get install docker.io python-pip curl
```

`docker.io` is a ubuntu 14.04 package refer to Docker getting started for
[instalations instructions]
(https://www.docker.io/gettingstarted/#h_installation).

To make easy the use of a container we use [fig](http://orchardup.github.io/fig/).

```
sudo pip install fig
```

Use the follow directory structure to keep MongoDB configuration, logs
and data.

```
mkdir -p ~/.containers/apollo/volumes/mongodb/log
mkdir -p ~/.containers/apollo/volumes/mongodb/lib
mkdir -p ~/.containers/apollo/volumes/mongodb/etc
curl https://raw.githubusercontent.com/wiliamsouza/docker-mongodb/master/volumes/etc/mongod.conf -o ~/.containers/apollo/volumes/mongodb/etc/mongod.conf
```

Start the container run:

```
fig up
```

The command should be executed inside root project directory.

If you need to access you MongoDB shell install the client tools.

```
sudo apt-get install mongodb-clients
```

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

From now you should set an environment variable `$TOKEN` it will be used through this doc:

```
export TOKEN=<token>
```

`<token>` is the one returned by `users/authenticate`.


Detail `user`:

```
curl -v --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request GET http://localhost:8000/users/jhon@doe.com
```
It will return `user` details:

```
{"name":"Jhon Doe","email":"jhon@doe.com","apikey":"0Yy00ZDRhLThmNmQt...","created":"2014-02-22T00:20:44.511-03:00","lastlogin":"0001-01-01T00:00:00Z"}
```

Package API
-----------

Uploading a `package`:

```
curl -vv --header "Authorization: Bearer $TOKEN" --request POST --form package=@package1.tgz --form metadata=@metadata1.json http://localhost:8000/tests/packages
```

The `data` directory contains some files to test.

It will return basic information about `package`:

```
{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}}
```

Listing `packages`:

```
curl -vv --header "Authorization: Bearer $TOKEN"  --header "Content-Type: application/json" -X GET http://localhost:8000/tests/packages
```

It will return a list of `packages`:

```
[{"filename":"package1.tgz","metadata":{"description":"Package1 ON/OFF test"}}]
```

Detailing `package`:

```
curl -vv --header "Authorization: Bearer $TOKEN"  --header "Content-Type: application/json" -X GET http://localhost:8000/tests/packages/package1.tgz
```

It will return `package` detail:

```
{"filename":"package1.tgz","metadata":{"version":0.1,"description":"Package1 ON/OFF test","install":"adb push dist/package1.jar /data/local/tmp/","run":"adb shell uiautomator runtest package1.jar -c com.github.wiliamsouza.package1.Package1Test"}}
```

Downloading `package`:

```
curl -vv --header "Authorization: Bearer $TOKEN"  -X GET http://localhost:8000/tests/packages/downloads/package1.tgz -o package.tgz
```

It will download `package`.

Organization API
----------------

Adding a new `organization`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request POST --data '{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}' http://localhost:8000/organizations
```

It will return:

```
{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}
```

List `organizations`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request GET http://localhost:8000/organizations
```

It will return:

```
[{"name":"janecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jane@doe.com"]},{"name":"jhoncorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}]
```

Detail `organization`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request GET http://localhost:8000/organizations/doecorp
```

It will return:

```
{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com","jane@doe.com"]}],"admins":["jhon@doe.com"]}
```

Update `organization`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request PUT --data '{"name":"doecorp","teams":[{"name":"Test","users":["jhon@doe.com"]}],"admins":["jane@doe.com"]}' http://localhost:8000/organizations/doecorp
```

Update is used to simulate deletion, change and addtion.

Delete:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request DELETE http://localhost:8000/organizations/doecorp
```

Device API
----------

Adding a new device:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request POST --data '{"codename":"a700","permission":{"organization":{"run":false,"result":false,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"}' http://localhost:8000/devices
```

It will return:

```
{"codename":"a700","permission":{"organization":{"run":false,"result":false,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"}
```

List `devices`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request GET http://localhost:8000/devices
```

It will return:

```
[{"codename":"a700","permission":{"organization":{"run":false,"result":false,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"}]
```

Detail `device`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request GET http://localhost:8000/devices/a700
```

It will return:

```
{"codename":"a700","permission":{"organization":{"run":false,"result":false,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"}
```

Update `device`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request PUT --data '{"codename":"a700","permission":{"organization":{"run":true,"result":true,"info":false},"team":{"run":false,"result":false,"info":false}},"owner":"jhon@doe.com","status":"","name":"Acer A700","vendor":"acer","manufacturer":"acer","type":"tablet","platform":"NVIDIA Tegra 3","cpu":"1.3 GHz quad-core Cortex A9","gpu":"416 MHz twelve-core Nvidia GeForce ULP","ram":"1GB","weight":"665 g (1.47 lb)","dimensions":"259x175x11 mm (10.20x6.89x0.43 in)","screenDimension":"257 mm (10.1 in)","resolution":"1920x1200","screenDensity":"224 PPI","internalStorage":"32GB","sdCard":"up to 32 GB","bluetooth":"yes","wifi":"802.11 b/g/n","mainCamera":"5MP","secondaryCamera":"2MP","power":"9800 mAh","peripherals":"accelerometer, gyroscope, proximity sensor, digital compass, GPS, magnometer, microphone"}' http://localhost:8000/devices/a700
```

Update is used to simulate deletion, change and addtion.

Delete:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request DELETE http://localhost:8000/devices/a700
```

Cicle API
---------

Adding a new `cicle`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request POST --data '{"name":"Test maguro","device":"maguro","packages":["bluetooth.jar","wifi.jar"]}' http://localhost:8000/tests
```

It will return:

```
{"id":"530d6474e169975bbb000001","name":"Test maguro","device":"maguro","packages":["bluetooth.jar","wifi.jar"]}
```

List `cicles`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request GET http://localhost:8000/tests
```

It will return:

```
[{"id":"530d6474e169975bbb000001","name":"Test maguro","device":"maguro","packages":["bluetooth.jar","wifi.jar"]}]
```

Detail `cicle`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request GET http://localhost:8000/tests/530d6474e169975bbb000001
```

It will return:

```
{"id":"530d6474e169975bbb000001","name":"Test maguro","device":"maguro","packages":["bluetooth.jar","wifi.jar"]}
```

Update `cicle`:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request PUT --data '{"name":"Test yakju","device":"yakju","packages":["bluetooth.jar","wifi.jar"]}' http://localhost:8000/tests/530d6474e169975bbb000001
```

Update is used to simulate deletion, change and addtion.

Delete:

```
curl -vv --header "Content-Type: application/json" --header "Authorization: Bearer $TOKEN" --request DELETE http://localhost:8000/tests/530d6474e169975bbb000001
```

[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/wiliamsouza/apollo/trend.png)](https://bitdeli.com/free "Bitdeli Badge")
