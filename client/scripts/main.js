import { connect } from './net.js';
import { Camera } from './camera.js';
import { render } from './renderer.js';
import { initInput, sendMouseUpdate } from './input.js';
import { getInterpolated } from './interpolation.js';

import './vec2.js';
import './utils.js';
import './interpolation.js';
import './net.js';

const canvas = /** @type {HTMLCanvasElement} */ (document.getElementById('canvas'));
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;
window.addEventListener('resize', () => {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
})


connect();

const camera = new Camera();

initInput(canvas, camera);

function gameLoop() {
    const snapshot = getInterpolated();
    if (snapshot) {
        sendMouseUpdate(canvas, camera);

        camera.update(snapshot, canvas);

        render(snapshot, camera, canvas);
    }
    requestAnimationFrame(gameLoop);
}
requestAnimationFrame(gameLoop);

