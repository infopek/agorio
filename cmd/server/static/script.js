let ws = null;

const canvas = document.getElementById("gameCanvas");
const ctx = canvas.getContext("2d");
const messagesDiv = document.getElementById("messages");

const GRID_CELL_SIZE = 30;

// TODO: movable camera
const camera = {
    x: 0,
    y: 0
};

let currPlayerID = null;

function connect() {
    ws = new WebSocket("ws://localhost:8080/ws");
    ws.onopen = () => console.log("connected");
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === "init") {  // should be sent only once when connected
            currPlayerID = data.player_id;
            return;
        }

        const me = data.players[currPlayerID];
        if (me) {
            clear();

            updateCamera(me.cells);

            draw(data);
        }
    };
    ws.onclose = () => setTimeout(connect, 1000);
    ws.onerror = (err) => console.error(err);
}


function draw(data) {
    drawGrid();
    drawPlayers(data);
}

function drawPlayers(data) {
    const players = data.players;

    const me = players[currPlayerID];   // currently unused
    console.log("world coords: ", me.cells[0].body.position);
    for (const [playerID, player] of Object.entries(players)) {
        for (const cell of player.cells) {
            const pos = cell.body.position;
            const screenCoords = worldToScreen(pos.x, pos.y);

            ctx.beginPath();
            ctx.arc(screenCoords.x, screenCoords.y, cell.body.mass, 0, 2 * Math.PI);
            ctx.strokeStyle = getPlayerColor(playerID);
            ctx.stroke();
        }
    }
}

function drawGrid() {
    // Camera rect
    const worldLeft   = camera.x - canvas.width / 2;
    const worldRight  = camera.x + canvas.width / 2;
    const worldTop    = camera.y - canvas.height / 2;
    const worldBottom = camera.y + canvas.height / 2;

    // Nearest grid line to avoid missing lines on edges
    const startX = Math.floor(worldLeft / GRID_CELL_SIZE) * GRID_CELL_SIZE;
    const endX = Math.ceil(worldRight / GRID_CELL_SIZE) * GRID_CELL_SIZE;
    const startY = Math.floor(worldTop / GRID_CELL_SIZE) * GRID_CELL_SIZE;
    const endY = Math.ceil(worldBottom / GRID_CELL_SIZE) * GRID_CELL_SIZE;

    ctx.save();
    ctx.strokeStyle = "#cccccc";
    ctx.lineWidth = 1;

    // Vertical lines
    for (let x = startX; x <= endX; x += GRID_CELL_SIZE) {
        const screenStart = worldToScreen(x, worldTop);
        const screenEnd = worldToScreen(x, worldBottom);

        ctx.beginPath();
        ctx.moveTo(screenStart.x, screenStart.y);
        ctx.lineTo(screenEnd.x, screenEnd.y);
        ctx.stroke();
    }

    // Horizontal lines
    for (let y = startY; y <= endY; y += GRID_CELL_SIZE) {
        const screenStart = worldToScreen(worldLeft, y);
        const screenEnd = worldToScreen(worldRight, y);

        ctx.beginPath();
        ctx.moveTo(screenStart.x, screenStart.y);
        ctx.lineTo(screenEnd.x, screenEnd.y);
        ctx.stroke();
    }

    ctx.restore();
}

function clear() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
}

function getPlayerColor(playerID) {
    switch (playerID) {
        case currPlayerID:
            return "green";
        default:
            return "red";
    }
}

function worldToScreen(worldX, worldY) {
    return {
        x: worldX - camera.x + canvas.width / 2,
        y: worldY - camera.y + canvas.height / 2
    };
}

function updateCamera(cells) {
    if (!cells.length) {
        return;
    }

    let sumX = 0;
    let sumY = 0;
    for (const cell of cells) {
        sumX += cell.body.position.x;
        sumY += cell.body.position.y;
    }

    camera.x = sumX / cells.length;
    camera.y = sumY / cells.length;
}

canvas.addEventListener("mousemove", (event) => {
    if (!ws || ws.readyState !== WebSocket.OPEN) {
        return;
    }

    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    ws.send(JSON.stringify({
        type: "move",
        x:    x,
        y:    y
    }));
});

connect();
