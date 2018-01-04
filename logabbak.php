<html>
<head>
<meta charset="UTF-8" />
<script>
var loginname = "";
// var serversocket = new WebSocket("ws://localhost:8081/send");
// var serversocket = "";  
var serversocketa = "";  
var serversocketb = "";  
var messagesent = true;
var messagequeue = [];
 
function utf8_to_b64( str ) {
  return window.btoa(unescape(encodeURIComponent( str )));
}

function b64_to_utf8( str ) {
  return decodeURIComponent(escape(window.atob( str )));
}

function connectws(serversocket, div){

    document.getElementById(div).innerHTML += "connectws   <br>";    
    serversocket = new WebSocket("ws://localhost:8080/send");

  serversocket.onopen = function() {
          // serversocket.send("{\"email\" : \"world\"}");
          // sendlogin() 
           // document.getElementById(div).innerHTML += "socket open <br>";    
            print("onopen "+div+" readyState "+serversocket.readyState, div );
   
  }
  serversocket.onclose = function() {
          // serversocket.send("{\"email\" : \"world\"}");
          // sendlogin() 
           document.getElementById(div).innerHTML += "socket closed <br>";    

  }
  // Write message on receive
  serversocket.onmessage = function(e) {
          document.getElementById(div).innerHTML += "Received: " + e.data + "<br>";
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


  function sendws(data, serversocket){
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

  function print (data, div){
    document.getElementById(div).innerHTML +=  data + "<br>";    
  }


  function senddata_ab() {
          var data = document.getElementById('sendtext').value;
          var timeStamp = Math.floor(Date.now()); 
          var uu = 'msg'+timeStamp+'.'+utf8_to_b64('a@a.com,b@a.com');
          data = '{ "type":"msg", "from":"a@a.com","to":"b@a.com","msg":"ola from b" ,"time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
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
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64(loginname+',system,makegroupadmin');
    var data = '{ "type":"system", "from":"'+loginname+'","to":"system","msg":"makegroupadmin,g123456" ,"cmd": "makegroupadmin", "groupname":"'+groupname+'", "member":"'+membername+'","time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
    sendws(data);
  }

  function getgrouplist() { 
    var timeStamp = (Date.now()); 
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64(loginname+',system,getgrouplist');
    var data = '{ "type":"system", "from":"'+loginname+'","to":"system","msg":"getgrouplist,g123456" ,"cmd": "makegroupadmin", "groupname":"a@a.com#g123456",  "time":'+timeStamp+'   ,"uuid":"'+uu+'"}';
    sendws(data);  
  }

  function getuserinfo() {   
  }

  function setuserinfo() {   
  }
  
  function uploadfile() {   
  }



  function sendlogin(aa, serversocket) {
     var timeStamp = (Date.now()); 
    var uu = 'msg'+timeStamp+'.'+utf8_to_b64(aa+',system,getgrouplist');
           data = '{"type":"login",  "login":"'+aa+'" ,"uuid":"'+uu+'"}';
          sendws(data, serversocket);
           
  } 
  


function generateUUID() {
    var d = new Date().getTime();
    var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = (d + Math.random()*16)%16 | 0;
        d = Math.floor(d/16);
        return (c=='x' ? r : (r&0x3|0x8)).toString(16);
    });
    return uuid;
};

 function testlogin(serversocket) {   
    
  // connectws();
  sendlogin('a@a.com', serversocket);
  // print (serversocket.readyState);
  }

  function testfunctions() {   
     
  // senddata_ab();
   print (serversocket.readyState);
 
  creategroup("a@a.com#grouptest", 'a@a.com');
  addtogroup("a@a.com#grouptest", 'a@a.com', 'b@a.com');

  for (var i = 0; i < 1000; i++) {
    addtogroup("a@a.com#grouptest", 'a@a.com', 'c'+i+'@a.com'); 
  }
  deletefromgroup("a@a.com#grouptest", 'a@a.com', 'member');
// makegroupadmin("a@a.com#grouptest", 'b@a.com');
  }

</script>
</head>
 
<body>

         
        <div style="width:50%;height: 100%;position: absolute;left: 0px;top:0px;background: rgb(255,200,200);">
        <input type="button" id="connect" value="connectws" onclick="connectws(serversocketa,'comms')"></input> 
      <!--   <input type="button" id="sendBtn" value="sendlogina" onclick="sendlogin('a@a.com')"></input> 
        <input type="button" id="sendBtn" value="sendloginb" onclick="sendlogin('b@a.com')"></input>
        <input type="button" id="sendBtn" value="senddata_ab" onclick="senddata_ab()"></input>
        <input type="button" id="sendBtn" value="senddata_ba" onclick="senddata_ba()"></input>
        <input type="button" id="sendBtn" value="sendgroup" onclick="sendgroup()"></input>
        <input type="button" id="sendBtn" value="creategroup" onclick="creategroup()"></input>
        <input type="button" id="sendBtn" value="addtogroup" onclick="addtogroup()"></input>
        <input type="button" id="sendBtn" value="removefromgroup" onclick="removefromgroup()"></input> -->
        <input type="button" id="sendBtn" value="testlogin" onclick="testlogin(serversocketa)"></input> 
        <input type="button" id="sendBtn" value="test" onclick="testfunctions()"></input>

        <div id='comms' style="width:100%;"></div>
        </div>
        <div style="width:50%;height: 100%;position: absolute;right: 0px;top:0px;background: rgb(255,255,200);">
        <input type="button" id="connect" value="connectws" onclick="connectws(serversocketb,'commsb')"></input> 
      <!--   <input type="button" id="sendBtn" value="sendlogina" onclick="sendlogin('a@a.com')"></input> 
        <input type="button" id="sendBtn" value="sendloginb" onclick="sendlogin('b@a.com')"></input>
        <input type="button" id="sendBtn" value="senddata_ab" onclick="senddata_ab()"></input>
        <input type="button" id="sendBtn" value="senddata_ba" onclick="senddata_ba()"></input>
        <input type="button" id="sendBtn" value="sendgroup" onclick="sendgroup()"></input>
        <input type="button" id="sendBtn" value="creategroup" onclick="creategroup()"></input>
        <input type="button" id="sendBtn" value="addtogroup" onclick="addtogroup()"></input>
        <input type="button" id="sendBtn" value="removefromgroup" onclick="removefromgroup()"></input> -->
        <input type="button" id="sendBtn" value="testlogin" onclick="testloginb(serversocketb)"></input> 
        <input type="button" id="sendBtn" value="test" onclick="testfunctionsb()"></input>

        <div id='commsb' style="width:100%;"></div>
        </div>
</body>
</html>