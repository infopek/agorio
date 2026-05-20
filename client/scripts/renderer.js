import { Vec2 } from './vec2.js'

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

    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Background
    drawGrid(ctx, camera, canvas);

    // Pellets

    // Viruses

    // Players
    for (let i = 0; i < snapshot.cells.length; i++) {
        const cell = snapshot.cells[i];
        const screenCoords = camera.worldToScreen(new Vec2(cell.x, cell.y), canvas);
        const screenRadius = camera.worldToScreenRadius(cell.radius);

        drawCell(ctx, screenCoords.x, screenCoords.y, screenRadius);
    }

    ctx.fillStyle = 'black';
    ctx.font = '20 monospace';
    ctx.fillText(`x: ${Math.round(camera.position.x)} y: ${Math.round(camera.position.y)}`, 10, 20);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./camera.js').Camera} camera
 * @param {HTMLCanvasElement} canvas
 */
function drawGrid(ctx, camera, canvas) {

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
