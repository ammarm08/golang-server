var app = require('express')();
var http = require('http').Server(app);
var io = require('socket.io')(http);

app.get('/', function(req, res){
  res.sendfile('index.html');
});

io.on('connection', function(socket){
  console.log('a user connected');

  socket.on('message', function(msg) {
    console.log('message received: ' + msg);
    io.emit('message', msg);
  })


  socket.on('disconnect', function() {
    io.emit('message', 'a user disconnected');
  })
});

http.listen(3000, function(){
  console.log('listening on *:3000');
});