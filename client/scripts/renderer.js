import { Vec2 } from './vec2.js';
import { Camera } from './camera.js';

const GRID_SIZE = 50;

/**
 * @param {import('./snapshot.js').Snapshot} snapshot
 * @param {import('./camera.js').Camera} camera
 * @param {HTMLCanvasElement} canvas
 */
export function render(snapshot, camera, canvas) {
    const ctx = canvas.getContext('2d');
    if (ctx === null) {
        return;
    }

    clear(ctx, canvas);

    // Background
    drawGrid(ctx, camera, canvas);

    // Pellets
    for (let i = 0; i < snapshot.pellets.length; i++) {
        const pellet = snapshot.pellets[i];
        const screenCoords = camera.worldToScreen(new Vec2(pellet.x, pellet.y), canvas);
        const screenRadius = camera.worldToScreenRadius(pellet.radius);

        drawPellet(ctx, screenCoords.x, screenCoords.y, screenRadius);
    }

    // Viruses
    for (let i = 0; i < snapshot.viruses.length; i++) {
        const virus = snapshot.viruses[i];
        const screenCoords = camera.worldToScreen(new Vec2(virus.x, virus.y), canvas);
        const screenRadius = camera.worldToScreenRadius(virus.radius);

        drawVirus(ctx, screenCoords.x, screenCoords.y, screenRadius);
    }

    // Players
    for (let i = 0; i < snapshot.cells.length; i++) {
        const cell = snapshot.cells[i];
        const screenCoords = camera.worldToScreen(new Vec2(cell.x, cell.y), canvas);
        const screenRadius = camera.worldToScreenRadius(cell.radius);

        drawCell(ctx, screenCoords.x, screenCoords.y, screenRadius);
    }

    displayWorldPosition(ctx, camera);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {HTMLCanvasElement} canvas
 */
function clear(ctx, canvas) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./camera.js').Camera} camera
 * @param {HTMLCanvasElement} canvas
 */
function drawGrid(ctx, camera, canvas) {
    // Camera rect
    const worldLeft = camera.position.x - canvas.width / 2 / camera.zoom;
    const worldRight = camera.position.x + canvas.width / 2 / camera.zoom;
    const worldTop = camera.position.y - canvas.height / 2 / camera.zoom;
    const worldBottom = camera.position.y + canvas.height / 2 / camera.zoom;

    // Set drawing cursor to top left
    const startX = Math.floor(worldLeft / GRID_SIZE) * GRID_SIZE;
    const startY = Math.floor(worldTop / GRID_SIZE) * GRID_SIZE;

    ctx.strokeStyle = '#222';
    ctx.lineWidth = 0.5;
    ctx.beginPath();
    for (let x = startX; x <= worldRight; x += GRID_SIZE) {
        const sx = (x - camera.position.x) * camera.zoom + canvas.width / 2;
        ctx.moveTo(sx, 0);
        ctx.lineTo(sx, canvas.height);
    }
    for (let y = startY; y <= worldBottom; y += GRID_SIZE) {
        const sy = (y - camera.position.y) * camera.zoom + canvas.height / 2;
        ctx.moveTo(0, sy);
        ctx.lineTo(canvas.width, sy);
    }
    ctx.stroke();
}

/** === FOR NOW, WE DRAW EVERYTHING AS CIRCLES XD === */

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 */
function drawPellet(ctx, x, y, r) {
    ctx.beginPath();
    ctx.arc(x, y, r, 0, 2 * Math.PI);
    ctx.stroke();
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 */
function drawVirus(ctx, x, y, r) {
    ctx.beginPath();
    ctx.arc(x, y, r, 0, 2 * Math.PI);
    ctx.stroke();
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 */
function drawCell(ctx, x, y, r) {
    ctx.beginPath();
    ctx.arc(x, y, r, 0, 2 * Math.PI);
    ctx.stroke();
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {Camera} camera
 */
function displayWorldPosition(ctx, camera) {
    ctx.fillStyle = 'black';
    ctx.font = '20 monospace';
    ctx.fillText(`x: ${Math.round(camera.position.x)} y: ${Math.round(camera.position.y)}`, 10, 20);
}


