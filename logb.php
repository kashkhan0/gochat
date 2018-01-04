<html><?php  
$login="a@a.com"; 
if(isset($_GET["login"]) )$login=$_GET["login"];  
$csrf = (isset($_GET["csrf"]) )? $_GET["csrf"] : "abdd";
$pass = (isset($_GET["pass"]) )? $_GET["pass"] : "abdd";
 
?>
<head>
<meta charset="UTF-8" />
<style type="text/css">
  body {
  font-family: Arial, "Helvetica Neue", Helvetica, sans-serif;
  font-size: 10px;
  font-style: normal;
  font-variant: normal;
  
}

</style>
<!-- <script src="ulid.js"></script> -->
<script>
var login = "<?php echo($login); ?>";
var csrf = "<?php echo($csrf); ?>";
var loginname = "";
// var serversocket = new WebSocket("ws://localhost:8081/send");
var serversocket = "";  
var messagesent = true;
var messagequeue = [];


  var ENCODING =  "0123456789abcdefghjkmnpqrstvwxyz"
  var ENCODING_LEN = ENCODING.length
  var RANDOM_LEN = 16

  function kencodeTime(time, len) {
    var mod
    var time
    var str = ""
    for (var x = len; x > 0; x--) {
      mod = time % ENCODING_LEN
      str = ENCODING.charAt(mod) + str
      time = (time - mod) / ENCODING_LEN
    }
    return str
  }
  function kencodeRandom(len) {
    var rand
    var str = ""
    for (var x = 0; x < len; x++) {
      // rand = Math.floor(ENCODING_LEN * prng())
      rand = Math.floor(ENCODING_LEN*crypto.getRandomValues(new Uint16Array(1))[0] / 0xFFFF)
      //rand = Math.floor(ENCODING_LEN*Math.random())
      str = ENCODING.charAt(rand) + str
    }
    return str
  }

  function ulid() {
      seedTime = Date.now();
    return kencodeTime(seedTime, 10) + kencodeRandom(16)
  }  

  // function ulids() {
  //     seedTime = Date.now();
  //   return kencodeTime(seedTime, 10) + kencodeRandom(6)
  // } 

  function ulids() {
      seedTime = Date.now();
    return ""+ seedTime + kencodeRandom(11)
  }
 
function utf8_to_b64( str ) {
  return window.btoa(unescape(encodeURIComponent( str )));
}

function b64_to_utf8( str ) {
  return decodeURIComponent(escape(window.atob( str )));
}


  function connectws(){
    document.getElementById('comms').innerHTML += "connectws   <br>";    
    serversocket = new WebSocket("ws://localhost:8080/send");

  serversocket.onopen = function() {
          // serversocket.send("{\"email\" : \"world\"}");
          // sendlogin() 
           document.getElementById('comms').innerHTML += "socket open <br>";    
            print("onopen serversocket.readyState "+serversocket.readyState );
   
  }
  serversocket.onclose = function() {
          // serversocket.send("{\"email\" : \"world\"}");
          // sendlogin() 
           document.getElementById('comms').innerHTML += "socket closed <br>";    

  }
  // Write message on receive
  serversocket.onmessage = function(e) {
          document.getElementById('comms').innerHTML += "Received: " + e.data + "<br>";
          // serversocket.send("received");
          var obj = JSON.parse(e.data);
           // document.getElementById('comms').innerHTML += "onmessage.selfack: " + obj.uuid  + "<br>";
          // if (obj.type =="system" && obj.msg =="loginok"){
          //   window.loginname = obj.to ;
          //   serversocket.send('{"type":"ackreceived","uuid":"'+ obj.uuid  +'" }');
          // }

          if (obj.type =="msg" ){
            serversocket.send('{"type":"msgackreceived","uuid":"'+ obj.uuid  +'" }');
            // document.getElementById('comms').innerHTML += "Received: " + e.data + "<br>";
          
          }


  };

}


  function sendws(data){
    // print("serversocket.readyState "+serversocket.readyState );
    if (serversocket.readyState != 1){
       print("socket not ready " );
       return;
    }
    serversocket.send(data);
    print("sendws sent: " + data );
  }

  function print (data){
    document.getElementById('comms').innerHTML +=  data + "<br>";    
  }

  function senddata_ab() {
          var data = document.getElementById('sendtext').value;
          var timeStamp = Math.floor(Date.now()); 
          var uu = 'msg'+timeStamp+'.'+utf8_to_b64('a@a.com,b@a.com');
          data = '{ "type":"msg", "from":"a@a.com","to":"b@a.com","msg":"ola from b","imageurl":"e2.jpg" ,"time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
          sendws(data);
  }
  function senddata_ba() {
          var data = document.getElementById('sendtext').value;
          var timeStamp = Math.floor(Date.now()); 
          var uu = 'msg'+timeStamp+'.'+utf8_to_b64('b@a.com,a@a.com');
          data = '{ "type":"msg", "from":"b@a.com","to":"a@a.com","msg":"ola from b" ,"time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
          sendws(data);
   } 
   function sendgroup() {
    if (loginname == "") {return}
          // var data = document.getElementById('sendtext').value;
          // var timeStamp = Math.floor(Date.now()); 
          // var uu = 'msg'+timeStamp+'.'+utf8_to_b64('b@a.com,group,b@a.com#');
          // data = '{ "type":"msggroup", "from":"b@a.com","to":"g123","msg":"ola from b" ,"time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
          // serversocket.send(data);
          // document.getElementById('comms').innerHTML += "senddata sent: " + data + "<br>";    
  } 


  function creategroup(groupname, username) {
    var timeStamp = (Date.now()); 
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64(username+',system,creategroup');
    var data = '{ "type":"system", "from":"'+username+'","to":"system","msg":"creategroup", "cmd": "creategroup", "groupname":"'+groupname+'" ,"time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
    sendws(data);
    }

  function addtogroup(groupname, username, membername) {
    // a@a.com#g123456 b@a.com
    var timeStamp = (Date.now()); 
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64(username+',system,addtogroup');
    var data = '{ "type":"system", "from":"'+username+'","to":"system","msg":"addtogroup" ,"cmd": "addtogroup", "groupname":"'+groupname+'", "membername":"'+membername+'", "membertype":"member" ,"time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
    sendws(data);
  }

  function deletefromgroup(groupname, username, membername) 
  {
    var timeStamp = (Date.now()); 
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64(username+',system,deletefromgroup');
    var data = '{ "type":"system", "from":"'+username+'","to":"system","msg":"creategroup,g123456" ,"cmd": "deletefromgroup", "groupname":"'+groupname+'", "membername":"'+membername+'","time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
    sendws(data);
  }

  function makegroupadmin(groupname, username, membername) 
  {   
    var timeStamp = (Date.now()); 
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64(username+',system,makegroupadmin');
    var data = '{ "type":"system", "from":"'+username+'","to":"system","msg":"makegroupadmin,g123456" ,"cmd": "makegroupadmin", "groupname":"'+groupname+'", "membername":"'+membername+'","time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
    sendws(data);
  }

  function getgrouplist() { 
    var timeStamp = (Date.now()); 
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64(username+',system,getgrouplist');
    var data = '{ "type":"system", "from":"'+username+'","to":"system","msg":"getgrouplist,g123456" ,"cmd": "makegroupadmin", "groupname":"a@a.com#g123456",  "time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
    sendws(data);  
  }

  function getuserinfo() {   
  }

  function setuserinfo() {   
  }
  
  function uploadfile() {   
  }



  function sendlogin(aa, acsrf) {
     var timeStamp = (Date.now()); 
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64(aa+',system,getgrouplist');
           data = '{"type":"login","login":"'+aa+'","csrf":"'+acsrf+'" ,"uuid":"'+uu+'"}';
          sendws(data);
           
  } 
  
function generateUUID() {
    return ulid();
};

function generateUUIDold() {
    var d = new Date().getTime();
    var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = (d + Math.random()*16)%16 | 0;
        d = Math.floor(d/16);
        return (c=='x' ? r : (r&0x3|0x8)).toString(16);
    });
    return uuid;
};

 function testbroadcast() {   
  var timeStamp = (Date.now()); 
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64('broadcast,system,getgrouplist');
           data = '{"type":"broadcast",  "to":"a@a.com" ,"from":"'+login+'","uuid":"'+uu+'"}';
          sendws(data);

 }
 function testmsg() {   
  var timeStamp = (Date.now()); 
    var uu = 'msg'+ulids();
           data = '{"type":"msg",  "to":"a@a.com" ,"msg":"message to a","imageurl":"e2.jpg" ,"from":"'+login+'","uuid":"'+uu+'"}';
          sendws(data);

 } 

 function testimg() {   
  var timeStamp = (Date.now()); 
  var uu = 'msg'+ulids();

  var img = '{ "name" :"a.jpg",  "contents" : "/9j/4AAQSkZJRgABAQIAHAAcAAD/2wBDACAWGBwYFCAcGhwkIiAmMFA0MCwsMGJGSjpQdGZ6eHJmcG6AkLicgIiuim5woNqirr7EztDOfJri8uDI8LjKzsb/2wBDASIkJDAqMF40NF7GhHCExsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsb/wAARCAAQABADASIAAhEBAxEB/8QAFgABAQEAAAAAAAAAAAAAAAAAAAIF/8QAHBAAAgEFAQAAAAAAAAAAAAAAAQIDAAQREkEx/8QAFQEBAQAAAAAAAAAAAAAAAAAAAgT/xAAbEQACAQUAAAAAAAAAAAAAAAAAEQMSIjGx0f/aAAwDAQACEQMRAD8Ax1t1C5kcKT5RrdSuY3DEe1Ec+q6uocczykk+y6ooQdx2hcypw042+H//2Q==" }';

           data = '{"type":"msg",  "to":"a@a.com" ,"msg":"message to a","imageurl":"e2.jpg" ,"from":"'+login+'","uuid":"'+uu+'"}';
          sendws(data);

 } 

 function testgroupmsg() {   
  var timeStamp = (Date.now()); 
    var uu = 'msg'+ulids();
           data = '{"type":"msg", "group":"a@a.com#grouptest" ,"msg":"message to group", "to":"group" ,"from":"'+login+'","imurl":["users/a@a.com/chatimages/imout1498808521802dwjnmv.jpg","users/a%40a.com/chatimages/imout1498832182153fq9e54.jpg"], "uuid":"'+uu+'"}';
          sendws(data);

 }
 function testlogin(username) {   
    
  // connectws();
  sendlogin(username, csrf);
   var timeStamp = Date.now(); 
    var uu = 'msg'+ulids();
  var data = '{ "type":"system", "from":"'+username+'","to":"system","msg":"getcachedmsg" ,"cmd": "getcachedmsg", "time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
 
  sendws(data)
  // print (serversocket.readyState);
  }

  function testfunctions() {   
     
  // senddata_ab();
   print (serversocket.readyState);
 
  creategroup("a@a.com#grouptest", 'a@a.com');
  addtogroup("a@a.com#grouptest", 'a@a.com', 'b@a.com');

  for (var i = 0; i < 4; i++) {
    addtogroup("a@a.com#grouptest", 'a@a.com', 'c'+i+'@a.com'); 
  }
  deletefromgroup("a@a.com#grouptest", 'a@a.com', 'c1@a.com');
  makegroupadmin("a@a.com#grouptest", 'a@a.com','b@a.com');
  }

</script>
</head>
 
<body>
        <!-- <textarea id="sendtext" type="text" > </textarea> -->
<?php echo $login." ".$csrf."<br>"; ?>
        <input type="button" id="connect" value="connectws" onclick="connectws()"></input> 
      <!--   <input type="button" id="sendBtn" value="sendlogina" onclick="sendlogin('a@a.com')"></input> 
        <input type="button" id="sendBtn" value="sendloginb" onclick="sendlogin('b@a.com')"></input>
        <input type="button" id="sendBtn" value="senddata_ab" onclick="senddata_ab()"></input>
        <input type="button" id="sendBtn" value="senddata_ba" onclick="senddata_ba()"></input>
        <input type="button" id="sendBtn" value="sendgroup" onclick="sendgroup()"></input>
        <input type="button" id="sendBtn" value="creategroup" onclick="creategroup()"></input>
        <input type="button" id="sendBtn" value="addtogroup" onclick="addtogroup()"></input>
        <input type="button" id="sendBtn" value="removefromgroup" onclick="removefromgroup()"></input> -->
        <input type="button" id="sendBtn" value="testlogin" onclick="testlogin('<?php echo $login;?>')"></input> 
        <input type="button" id="sendBtn" value="test" onclick="testfunctions()"></input>
<input type="button" id="sendBtn" value="testbroadcast" onclick="testbroadcast()"></input>
<input type="button" id="sendBtn" value="testmsg" onclick="testmsg()"></input>
<input type="button" id="sendBtn" value="testimgg" onclick="testimg()"></input>
<input type="button" id="sendBtn" value="testgroupmsg" onclick="testgroupmsg()"></input>

        <div id='comms' style="width:100%;height:80%;background: #ffa;overflow: auto;"></div>
</body>
</html>