<html>
<head>
<meta charset="UTF-8" />
<script>
var loginname = "";
var connection = new WebSocket('ws://localhost:8080/send');
  connection.onopen = function () {
      connection.send('Ping'); // Send the message 'Ping' to the server
    };

    // Log errors
    connection.onerror = function (error) {
      print('WebSocket Error ' + error);
    };

    // Log messages from the server
    connection.onmessage = function (e) {
      print('Server: ' + e.data);
    };


      function print (data){
    document.getElementById('comms').innerHTML +=  data + "<br>";    
  }


 // Sending String
connection.send('your message');

// Sending canvas ImageData as ArrayBuffer
var img = canvas_context.getImageData(0, 0, 400, 320);
var binary = new Uint8Array(img.data.length);
for (var i = 0; i < img.data.length; i++) {
  binary[i] = img.data[i];
}
connection.send(binary.buffer);

// Sending file as Blob
var file = document.querySelector('input[type="file"]').files[0];
connection.send(file);
</script>


</head>
<body>
        <textarea id="sendtext" type="text" > </textarea>/>
        <input type="button" id="sendBtn" value="testlogin" onclick="testlogin()"></input>
        <input type="button" id="sendBtn" value="test" onclick="testfunctions()"></input>

        <div id='comms' style="width:1200px;"></div>
</body>