import { sendJoin, sendMove } from './net.js';
import { Camera } from './camera.js';

const menu = document.getElementById('menu');
const nameInput = document.getElementById('name-input');
const playBtn = document.getElementById('play-btn');

let mouseX = 0;
let mouseY = 0;

/**
 * @param {HTMLCanvasElement} canvas
 * @param {Camera} camera
 */
export function sendMouseUpdate(canvas, camera) {
    const worldX = (mouseX - canvas.width / 2) / camera.zoom + camera.position.x;
    const worldY = (mouseY - canvas.height / 2) / camera.zoom + camera.position.y;
    sendMove(worldX, worldY);
}

/**
 * @param {HTMLCanvasElement} canvas
 * @param {Camera} camera
 */
export function initInput(canvas, camera) {
    if (!menu) {
        throw new Error('missing #menu');
    }
    if (!nameInput) {
        throw new Error('missing #name-input');
    }
    if (!playBtn) {
        throw new Error('missing #play-btn');
    }

    canvas.addEventListener('mousemove', (e) => {
         mouseX = e.clientX;
         mouseY = e.clientY;
    });

    playBtn.addEventListener('click', () => {
        // @ts-ignore
        const name = nameInput.value.trim() || 'Player';
        sendJoin(name);
        menu.classList.add('hidden');
    });
}

