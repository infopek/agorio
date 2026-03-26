let ws = null;

const canvas = document.getElementById('gameCanvas');
const ctx = canvas.getContext("2d");
const messagesDiv = document.getElementById('messages');

function connect() {
    ws = new WebSocket('ws://localhost:8080/ws');
    ws.onopen = () => console.log('Connected');
    ws.onmessage = (event) => {
        messagesDiv.insertAdjacentHTML('beforeend', `<p>${event.data}</p>`);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
    };
    ws.onclose = () => setTimeout(connect, 1000);
    ws.onerror = (err) => console.error(err);
}

function sendMessage() {
    const input = document.getElementById('messageInput');
    const msg = input.value;
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(msg);
    }
    input.value = '';
}

canvas.addEventListener('mousemove', (event) => {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        return;
    }

    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    ws.send(JSON.stringify(
        {
            type: 'move',
            x,
            y
        }
    ));
});

connect();
