import { sendJoin, sendMove } from './net.js';
import { Camera } from './camera.js';

const menu = document.getElementById('menu');
const nameInput = document.getElementById('name-input');
const playBtn = document.getElementById('play-btn');

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

    playBtn.addEventListener('click', () => {
        // @ts-ignore
        const name = nameInput.value.trim() || 'Player';
        sendJoin(name);
        menu.classList.add('hidden');

        canvas.addEventListener('mousemove', (e) => {
            const worldX = (e.clientX - canvas.width / 2) / camera.zoom + camera.position.x;
            const worldY = (e.clientY - canvas.height / 2) / camera.zoom + camera.position.y;
            sendMove(worldX, worldY);
        });
    });
}

