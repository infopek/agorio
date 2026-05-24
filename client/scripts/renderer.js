import { Vec2 } from './vec2.js';
import { Camera } from './camera.js';

const GRID_SIZE = 50;

const WORLD_WIDTH = 5000;
const WORLD_HEIGHT = 5000;

/** render
 *
 * Main render function that renders a snapshot
 *
 * @param {import('./interpolation.js').Snapshot} snapshot
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
    drawWorldBorder(ctx, camera, canvas);

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
        displayMass(ctx, screenCoords.x, screenCoords.y, screenRadius, cell.mass);
    }

    displayWorldPosition(ctx, camera);
}

/** clear
 *
 * Clears the screen
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {HTMLCanvasElement} canvas
 */
function clear(ctx, canvas) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
}

/** drawWorldBorder
 *
 * Draws the boundary of the world (hardcoded 5000x5000 constants)
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./camera.js').Camera} camera
 * @param {HTMLCanvasElement} canvas
 */
function drawWorldBorder(ctx, camera, canvas) {
    const topLeft = camera.worldToScreen(new Vec2(0, 0), canvas);
    const bottomRight = camera.worldToScreen(new Vec2(WORLD_WIDTH, WORLD_HEIGHT), canvas);

    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        ctx.strokeStyle = '#ff0000';
        ctx.lineWidth = 3;
        ctx.strokeRect(
            topLeft.x,
            topLeft.y,
            bottomRight.x - topLeft.x,
            bottomRight.y - topLeft.y,
        );
    });
}

/** drawGrid
 *
 * Draws a ruler-like grid on the background
 *
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

    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        ctx.strokeStyle = '#bbbbbb';
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
    });
}

/** === FOR NOW, WE DRAW EVERYTHING AS CIRCLES XD === **/

/** drawPellet
 *
 * Draws a pellet in the world
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 */
function drawPellet(ctx, x, y, r) {
    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        ctx.beginPath();
        ctx.strokeStyle = '#666666';
        ctx.arc(x, y, r, 0, 2 * Math.PI);
        ctx.stroke();
    });
}

/** drawVirus
 *
 * Draws a virus in the world
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 */
function drawVirus(ctx, x, y, r) {
    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        ctx.beginPath();
        ctx.arc(x, y, r, 0, 2 * Math.PI);
        ctx.stroke();
    });
}

/** drawCell
 *
 * Draws a cell in the world
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 */
function drawCell(ctx, x, y, r) {
    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        ctx.beginPath();
        ctx.arc(x, y, r, 0, 2 * Math.PI);
        ctx.stroke();
    });
}


/** === DEBUG MOSTLY (YET) === **/

/** displayWorldPosition
 *
 * Displays the (x, y) coords of the player ('s camera) in world space
 *  in the top left of the screen
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {Camera} camera
 */
function displayWorldPosition(ctx, camera) {
    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        ctx.fillStyle = 'black';
        ctx.font = '20 monospace';
        ctx.fillText(`x: ${Math.round(camera.position.x)} y: ${Math.round(camera.position.y)}`, 10, 20);
    });
}

/** displayMass
 *
 * Displays the mass of the cell in the middle of it
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} screenX
 * @param {number} screenY
 * @param {number} screenRadius
 * @param {number} mass
 */
function displayMass(ctx, screenX, screenY, screenRadius, mass) {
    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        ctx.fillStyle = 'black';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.font = `${Math.max(12, screenRadius * 0.4)}px sans-serif`;
        ctx.fillText(`${Math.round(mass)}`, screenX, screenY);
    });
}

/** === UTILS === **/

/** draw
 *
 * Middleware for stateless draw calls
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {function} fn
 */
function draw(ctx, fn) {
    ctx.save();
    fn(ctx);
    ctx.restore();
}
