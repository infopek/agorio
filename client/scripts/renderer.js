import { Vec2 } from './vec2.js';
import { Camera } from './camera.js';
import { Config } from './config.js';

const OUTLINE_DARKEN_AMOUNT = 0.3;
const SHAPE_OUTLINE_WIDTH = 3;
const GRID_LINE_WIDTH = 0.5;

const VIRUS_SPIKE_COUNT = 30;
const VIRUS_SPIKE_DEPTH_SCALE = 0.04;
const VIRUS_ROTATION_DIRECTION = 1; // clockwise
const VIRUS_ROTATION_SPEED = 0.3;
const VIRUS_MASS_ROTATION_FACTOR = 0.7;
const VIRUS_COUNTER_FONT_SCALE = 0.6;
const VIRUS_COUNTER_OUTLINE_WIDTH = 2;

const CELL_TEXT_Y_OFFSET_SCALE = 0.35;
const CELL_NAME_MIN_FONT_SIZE = 13;
const CELL_NAME_RADIUS_SCALE = 0.38;
const CELL_NAME_LENGTH_TARGET = 7;
const CELL_TEXT_OUTLINE_SCALE = 0.11;
const CELL_MASS_MIN_FONT_SIZE = 8;
const CELL_MASS_RADIUS_SCALE = 0.2;

const TEXT_FONT_FAMILY = 'Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif';
const DEBUG_POSITION_FONT_SIZE = 13;
const DEBUG_POSITION_X = 12;
const DEBUG_POSITION_Y = 12;
const DEBUG_POSITION_PADDING_X = 10;
const DEBUG_POSITION_PADDING_Y = 7;
const DEBUG_POSITION_RADIUS = 6;

const WOBBLY_CIRCLE_POINT_COUNT = 60;
const WOBBLY_CIRCLE_WAVE_COUNT = 13;
const WOBBLY_CIRCLE_SIN_SPEED = 6;
const WOBBLY_CIRCLE_COS_SPEED = 5;
const WOBBLY_CIRCLE_AMPLITUDE = 0.0015;

/** render
 *
 * Main render function that renders a snapshot
 *
 * @param {import('./interpolation.js').Snapshot} snapshot
 * @param {import('./camera.js').Camera} camera
 * @param {HTMLCanvasElement} canvas
 * @param {CanvasRenderingContext2D} ctx
 */
export function render(snapshot, camera, canvas, ctx) {
    const time = performance.now() / 1000.0;

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

        drawEject(ctx, screenCoords.x, screenCoords.y, screenRadius, eject.color, time);
    }

    // Viruses
    for (let i = 0; i < snapshot.viruses.length; i++) {
        const virus = snapshot.viruses[i];
        const screenCoords = camera.worldToScreen(new Vec2(virus.x, virus.y), canvas);
        const screenRadius = camera.worldToScreenRadius(virus.radius);

        drawVirus(ctx, screenCoords.x, screenCoords.y, screenRadius, virus.mass, virus.color,
            darken(virus.color, OUTLINE_DARKEN_AMOUNT), time);

        const remaining = Config.virusFeedToShoot - virus.fed_count;
        if (remaining < Config.virusFeedToShoot) {
            drawOutlinedText(ctx, String(remaining), screenCoords.x, screenCoords.y, `bold ${screenRadius * VIRUS_COUNTER_FONT_SCALE}px sans-serif`, VIRUS_COUNTER_OUTLINE_WIDTH);
        }
    }

    // Players
    const cells = snapshot.cells.sort((a, b) => a.radius - b.radius); // consistent rendering, larger is in foreground
    for (let i = 0; i < cells.length; i++) {
        const cell = cells[i];
        const screenCoords = camera.worldToScreen(new Vec2(cell.x, cell.y), canvas);
        const screenRadius = camera.worldToScreenRadius(cell.radius);
        const offset = screenRadius * CELL_TEXT_Y_OFFSET_SCALE;

        drawCell(ctx, screenCoords.x, screenCoords.y, screenRadius, cell.color, time);

        displayName(ctx, screenCoords.x, screenCoords.y, screenRadius, cell.owner_name);
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
    const bottomRight = camera.worldToScreen(new Vec2(Config.worldWidth, Config.worldHeight), canvas);

    draw(ctx, ctx => {
        ctx.strokeStyle = '#ff0000';
        ctx.lineWidth = SHAPE_OUTLINE_WIDTH;
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
    const startX = Math.floor(worldLeft / Config.gridSize) * Config.gridSize;
    const startY = Math.floor(worldTop / Config.gridSize) * Config.gridSize;

    draw(ctx, ctx => {
        ctx.strokeStyle = '#bbbbbb';
        ctx.lineWidth = GRID_LINE_WIDTH;
        ctx.beginPath();
        for (let x = startX; x <= worldRight; x += Config.gridSize) {
            const sx = (x - camera.position.x) * camera.zoom + canvas.width / 2;
            ctx.moveTo(sx, 0);
            ctx.lineTo(sx, canvas.height);
        }
        for (let y = startY; y <= worldBottom; y += Config.gridSize) {
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
    draw(ctx, ctx => {
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
    draw(ctx, ctx => {
        const spikeDepth = r * VIRUS_SPIKE_DEPTH_SCALE;

        const rotation = time * VIRUS_ROTATION_SPEED * VIRUS_ROTATION_DIRECTION
            + mass * VIRUS_MASS_ROTATION_FACTOR;

        ctx.beginPath();
        for (let i = 0; i <= VIRUS_SPIKE_COUNT * 2; i++) {
            const angle = (i / (VIRUS_SPIKE_COUNT * 2)) * Math.PI * 2 + rotation;
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
        ctx.lineWidth = SHAPE_OUTLINE_WIDTH;
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
 * @param {number} time
 */
function drawCell(ctx, x, y, r, color, time) {
    drawWobblyCircle(ctx, x, y, r, color, darken(color, OUTLINE_DARKEN_AMOUNT), time)
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
 * @param {number} time
 */
function drawEject(ctx, x, y, r, color, time) {
    drawWobblyCircle(ctx, x, y, r, color, darken(color, OUTLINE_DARKEN_AMOUNT), time)
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
    draw(ctx, ctx => {
        const text = `x ${Math.round(camera.position.x)}   y ${Math.round(camera.position.y)}`;
        ctx.font = `600 ${DEBUG_POSITION_FONT_SIZE}px ui-monospace, SFMono-Regular, Menlo, Consolas, monospace`;
        const width = Math.ceil(ctx.measureText(text).width + DEBUG_POSITION_PADDING_X * 2);
        const height = DEBUG_POSITION_FONT_SIZE + DEBUG_POSITION_PADDING_Y * 2;

        ctx.fillStyle = 'rgba(255, 255, 255, 0.82)';
        ctx.strokeStyle = 'rgba(23, 32, 51, 0.12)';
        ctx.lineWidth = 1;
        roundedRect(ctx, DEBUG_POSITION_X, DEBUG_POSITION_Y, width, height, DEBUG_POSITION_RADIUS);
        ctx.fill();
        ctx.stroke();

        ctx.fillStyle = 'rgba(23, 32, 51, 0.72)';
        ctx.textAlign = 'left';
        ctx.textBaseline = 'middle';
        ctx.fillText(text, DEBUG_POSITION_X + DEBUG_POSITION_PADDING_X, DEBUG_POSITION_Y + height / 2);
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
 */
function displayName(ctx, x, y, r, name) {
    const baseSize = Math.max(CELL_NAME_MIN_FONT_SIZE, r * CELL_NAME_RADIUS_SCALE);
    const scale = name.length === 0 ? 1.0 : Math.min(1.0, CELL_NAME_LENGTH_TARGET / name.length);
    const size = baseSize * scale;
    drawOutlinedText(ctx, name, x, y, `800 ${size}px ${TEXT_FONT_FAMILY}`, Math.max(1.0, size * CELL_TEXT_OUTLINE_SCALE));
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
    const size = Math.max(CELL_MASS_MIN_FONT_SIZE, r * CELL_MASS_RADIUS_SCALE);
    drawOutlinedText(ctx, String(Math.round(mass)), x, y, `700 ${size}px ${TEXT_FONT_FAMILY}`, Math.max(1.0, size * CELL_TEXT_OUTLINE_SCALE));
}

/** === UTILS === **/

/** draw
 *
 * Middleware for stateless draw calls
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {(ctx: CanvasRenderingContext2D) => void} fn
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
    draw(ctx, ctx => {
        ctx.font = font;
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.lineWidth = lineWidth;
        ctx.lineJoin = 'round';
        ctx.miterLimit = 2;
        ctx.strokeStyle = 'rgba(20, 25, 32, 0.78)';
        ctx.strokeText(text, x, y);
        ctx.fillStyle = 'rgba(255, 255, 255, 0.96)';
        ctx.fillText(text, x, y);
    });
}


/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} width
 * @param {number} height
 * @param {number} radius
 */
function roundedRect(ctx, x, y, width, height, radius) {
    const r = Math.min(radius, width / 2, height / 2);
    ctx.beginPath();
    ctx.moveTo(x + r, y);
    ctx.lineTo(x + width - r, y);
    ctx.quadraticCurveTo(x + width, y, x + width, y + r);
    ctx.lineTo(x + width, y + height - r);
    ctx.quadraticCurveTo(x + width, y + height, x + width - r, y + height);
    ctx.lineTo(x + r, y + height);
    ctx.quadraticCurveTo(x, y + height, x, y + height - r);
    ctx.lineTo(x, y + r);
    ctx.quadraticCurveTo(x, y, x + r, y);
    ctx.closePath();
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
    draw(ctx, ctx => {
        ctx.beginPath();
        for (let i = 0; i <= WOBBLY_CIRCLE_POINT_COUNT; i++) {
            const angle = (i / WOBBLY_CIRCLE_POINT_COUNT) * Math.PI * 2.0;
            const wobble = Math.sin(angle * WOBBLY_CIRCLE_WAVE_COUNT + time * WOBBLY_CIRCLE_SIN_SPEED) * r * WOBBLY_CIRCLE_AMPLITUDE
                + Math.cos(angle * WOBBLY_CIRCLE_WAVE_COUNT + time * WOBBLY_CIRCLE_COS_SPEED) * r * WOBBLY_CIRCLE_AMPLITUDE;
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
        ctx.lineWidth = SHAPE_OUTLINE_WIDTH;
        ctx.stroke();
    });
}
