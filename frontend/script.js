const serverUrl = 'ws://localhost:8080/ws/chat';  // Update with your WebSocket server URL
let socket;
let roomName = '';

function joinRoom() {
    roomName = document.getElementById('room-input').value;
    socket = new WebSocket(serverUrl);

    socket.onopen = function(e) {
        console.log("Connection established");
        socket.send(JSON.stringify({ type: 'join_room', room: roomName }));
    };

    socket.onmessage = function(event) {
        const chatBox = document.getElementById('chat-box');
        const data = JSON.parse(event.data);
        chatBox.innerHTML += `<div>${data.content}</div>`;
    };

    socket.onclose = function(event) {
        if (event.wasClean) {
            console.log(`Connection closed cleanly, code=${event.code}`);
        } else {
            console.error('Connection died');
        }
    };

    socket.onerror = function(error) {
        console.error(`[error] ${error.message}`);
    };
}

function leaveRoom() {
    if (socket) {
        socket.send(JSON.stringify({ type: 'leave_room', room: roomName }));
        socket.close();
        roomName = '';
    }
}

function sendMessage() {
    const messageInput = document.getElementById('message-input');
    const message = messageInput.value;
    socket.send(JSON.stringify({ type: 'chat_message', room: roomName, content: message }));
    messageInput.value = '';
}
