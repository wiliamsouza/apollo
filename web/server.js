var fs = require('fs');

var express = require('express');
var app = express();

app.use('/', express.static(__dirname));
app.use(express.bodyParser());

app.listen(3000);
