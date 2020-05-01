var username;
var ws = null;

function setUsername() {
    username = document.getElementById("username").value;
    document.getElementById("isUser").innerHTML = username;
    ws.send(username);
}

function send(data) {
    body = {
        username: username,
        message: data
    }
    ws.send(JSON.stringify(body));
}

function connect() {
    var url = "ws://localhost:8080/websocket";
    ws = new WebSocket(url);
    ws.addEventListener("open", (ev) => {
        console.log("connected");
    });
    ws.addEventListener("message", (ev) => {
        var data = JSON.parse(ev.data);
        document.getElementById("chat").innerHTML += `<div style="background-color:white;padding:1em;margin:1em;"><small>${data.username}:</small><br>${data.message}</div>`
    });
}