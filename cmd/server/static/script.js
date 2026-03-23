// WebSocket connection and game logic
let ws = null;
const canvas = document.getElementById('gameCanvas');
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

// Add mouse move handler once canvas is ready
canvas.addEventListener('mousemove', (event) => {
    if (!ws || ws.readyState !== WebSocket.OPEN) return;
    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    ws.send(JSON.stringify({ type: 'move', x, y }));
});

// Start the connection
connect();
