import { Vec2 } from './vec2.js';
import { Camera } from './camera.js';

const GRID_SIZE = 30;

const WORLD_WIDTH = 15000;
const WORLD_HEIGHT = 15000;

const VIRUS_FEED_TO_SHOOT = 7;

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

        drawPellet(ctx, screenCoords.x, screenCoords.y, screenRadius, pellet.color);
    }

    // Eject
    for (let i = 0; i < snapshot.ejects.length; i++) {
        const eject = snapshot.ejects[i];
        const screenCoords = camera.worldToScreen(new Vec2(eject.x, eject.y), canvas);
        const screenRadius = camera.worldToScreenRadius(eject.radius);

        drawEject(ctx, screenCoords.x, screenCoords.y, screenRadius, eject.color);
    }

    // Viruses
    for (let i = 0; i < snapshot.viruses.length; i++) {
        const virus = snapshot.viruses[i];
        const screenCoords = camera.worldToScreen(new Vec2(virus.x, virus.y), canvas);
        const screenRadius = camera.worldToScreenRadius(virus.radius);

        drawVirus(ctx, screenCoords.x, screenCoords.y, screenRadius, virus.mass, virus.color,
            darken(virus.color, 0.3), performance.now() / 1000.0);

        const remaining = VIRUS_FEED_TO_SHOOT - virus.fed_count;
        if (remaining < VIRUS_FEED_TO_SHOOT) {
            drawOutlinedText(ctx, String(remaining), screenCoords.x, screenCoords.y, `bold ${screenRadius * 0.6}px sans-serif`, 2);
        }
    }

    // Players
    for (let i = 0; i < snapshot.cells.length; i++) {
        const cell = snapshot.cells[i];
        const screenCoords = camera.worldToScreen(new Vec2(cell.x, cell.y), canvas);
        const screenRadius = camera.worldToScreenRadius(cell.radius);
        const offset = screenRadius * 0.35;

        drawCell(ctx, screenCoords.x, screenCoords.y, screenRadius, cell.color);

        displayName(ctx, screenCoords.x, screenCoords.y, screenRadius, cell.owner_name, camera.zoom);
        displayMass(ctx, screenCoords.x, screenCoords.y + offset, screenRadius, cell.mass);
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
 * @param {string} color
 */
function drawPellet(ctx, x, y, r, color) {
    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        ctx.beginPath();
        ctx.arc(x, y, r, 0, 2 * Math.PI);
        ctx.fillStyle = color;
        ctx.fill();
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
 * @param {number} mass
 * @param {string} color
 * @param {string} outlineColor
 * @param {number} time
 */
function drawVirus(ctx, x, y, r, mass, color, outlineColor, time) {
    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        const spikes = 30;
        const spikeDepth = r * 0.04;
        const rotationDirection = 1; // clockwise
        const rotationSpeed = 0.3;
        const massRotationFactor = 0.7;

        const rotation = time * rotationSpeed * rotationDirection + mass * massRotationFactor;

        ctx.beginPath();
        for (let i = 0; i <= spikes * 2; i++) {
            const angle = (i / (spikes * 2)) * Math.PI * 2 + rotation;
            const isSpike = i % 2 == 0;
            const dist = isSpike ? r + spikeDepth : r - spikeDepth;
            const px = x + Math.cos(angle) * dist;
            const py = y + Math.sin(angle) * dist;

            if (i === 0) {
                ctx.moveTo(px, py);
            } else {
                ctx.lineTo(px, py);
            }
        }
        ctx.closePath();

        ctx.fillStyle = color;
        ctx.fill();

        ctx.strokeStyle = outlineColor;
        ctx.lineWidth = 3;
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
 * @param {string} color
 */
function drawCell(ctx, x, y, r, color) {
    drawWobblyCircle(ctx, x, y, r, color, darken(color, 0.3), performance.now() / 1000.0)
}

/** drawEject
 *
 * Draws an ejected mass in the world
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 * @param {string} color
 */
function drawEject(ctx, x, y, r, color) {
    drawWobblyCircle(ctx, x, y, r, color, darken(color, 0.3), performance.now() / 1000.0)
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

/** displayName
 *
 * Displays the name of the player on the cell
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 * @param {string} name
 * @param {number} zoom
 */
function displayName(ctx, x, y, r, name, zoom) {
    const baseSize = Math.max(13, r * 0.38);
    const scale = Math.min(1.0, 7.0 / name.length);
    const size = baseSize * scale;
    drawOutlinedText(ctx, name, x, y, `bold ${size}px sans-serif`, Math.max(1.0, size * 0.15));
}

/** displayMass
 *
 * Displays the mass of the cell
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 * @param {number} mass
 */
function displayMass(ctx, x, y, r, mass) {
    const size = Math.max(8, r * 0.2);
    drawOutlinedText(ctx, String(Math.round(mass)), x, y, `${size}px sans-serif`, Math.max(1.0, size * 0.15));
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

/**
 * darken
 *
 * Returns a slightly darker color of the one given
 *
 * @param {string} hex
 * @param {number} amount
 */
function darken(hex, amount) {
    const r = Math.max(0, parseInt(hex.slice(1, 3), 16) * (1 - amount));
    const g = Math.max(0, parseInt(hex.slice(3, 5), 16) * (1 - amount));
    const b = Math.max(0, parseInt(hex.slice(5, 7), 16) * (1 - amount));
    return `rgb(${Math.round(r)},${Math.round(g)},${Math.round(b)})`;
}

/** drawOutlinedText
 *
 * Draws the text in white with a black outline to the specified position
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {string} text
 * @param {number} x
 * @param {number} y
 * @param {string} font
 * @param {number} lineWidth
 */
function drawOutlinedText(ctx, text, x, y, font, lineWidth) {
    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        ctx.font = font;
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.lineWidth = lineWidth;
        ctx.lineJoin = 'round';
        ctx.miterLimit = 2;
        ctx.strokeStyle = 'black';
        ctx.strokeText(text, x, y);
        ctx.fillStyle = 'white';
        ctx.fillText(text, x, y);
    });
}

/** drawWobblyCircle
 *
 * Draws a very cute circle with outline that wobbles (some math magic animation)
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} r
 * @param {string} color
 * @param {string} outlineColor
 * @param {number} time
 */
function drawWobblyCircle(ctx, x, y, r, color, outlineColor, time) {
    draw(ctx, (/**@type {CanvasRenderingContext2D} */ctx) => {
        const points = 60;

        ctx.beginPath();
        for (let i = 0; i <= points; i++) {
            const angle = (i / points) * Math.PI * 2.0;
            const wobble = Math.sin(angle * 13.0 + time * 6.0) * r * 0.0015
                + Math.cos(angle * 13.0 + time * 5.0) * r * 0.0015;
            const px = x + Math.cos(angle) * (r + wobble);
            const py = y + Math.sin(angle) * (r + wobble);

            if (i === 0) {
                ctx.moveTo(px, py);
            } else {
                ctx.lineTo(px, py);
            }

        }
        ctx.closePath();
        ctx.fillStyle = color;
        ctx.fill();
        ctx.strokeStyle = outlineColor;
        ctx.lineWidth = 3;
        ctx.stroke();
    });
}
