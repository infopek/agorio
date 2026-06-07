import { sendJoin, sendMove, sendSplit, sendFeed } from './net.js';
import { Camera } from './camera.js';
import { Config } from './config.js';
import { getState, setState, ClientStatus, subscribe } from './state.js';

const menu = document.getElementById('menu');
const messageEl = document.getElementById('menu-message');
const nameInput = /** @type {HTMLInputElement | null} */ (document.getElementById('name-input'));
const playBtn = /** @type {HTMLButtonElement | null} */ (document.getElementById('play-btn'));

let mouseX = 0;
let mouseY = 0;

/**
 * @param {HTMLCanvasElement} canvas
 * @param {Camera} camera
 */
export function sendMouseUpdate(canvas, camera) {
    if (getState().status !== ClientStatus.PLAYING) {
        return;
    }

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
    if (!messageEl) {
        throw new Error('missing #menu-message');
    }
    if (!nameInput) {
        throw new Error('missing #name-input');
    }
    if (!playBtn) {
        throw new Error('missing #play-btn');
    }

    subscribe((state) => {
        menu.classList.toggle('hidden', state.status === ClientStatus.PLAYING);
        messageEl.textContent = state.message;

        playBtn.textContent = getPlayButtonLabel(state.status);
        playBtn.disabled = state.status === ClientStatus.CONNECTING;

        nameInput.disabled =
            state.status === ClientStatus.CONNECTING
            || state.status === ClientStatus.DISCONNECTED
            || state.status === ClientStatus.PAUSED;
    })

    canvas.addEventListener('mousemove', (e) => {
        mouseX = e.clientX;
        mouseY = e.clientY;
    });

    canvas.addEventListener('wheel', (e) => {
        e.preventDefault();
        camera.zoomBy(e.deltaY > 0.0 ? Config.wheelZoomOutScale : Config.wheelZoomInScale);
    });

    document.addEventListener('keydown', (e) => {
        const state = getState();
        if (e.code === 'Escape') {
            if (state.status === ClientStatus.PLAYING) {
                setState({ status: ClientStatus.PAUSED, message: 'Paused' });
            } else if (state.status === ClientStatus.PAUSED) {
                setState({ status: ClientStatus.PLAYING, message: '' });
            }
            return;
        }

        if (isTypingTarget(e.target)) {
            if (e.code === 'Enter') {
                handlePlayClick();
            }
            return;
        }
        if (state.status !== ClientStatus.PLAYING) {
            return; // below inputs only valid when playing
        }

        if (e.code === 'Space') {
            sendSplit();
        }
        if (e.code === 'KeyW') {
            sendFeed();
        }
    });

    playBtn.addEventListener('click', handlePlayClick);
}

function handlePlayClick() {
    const state = getState();

    if (state.status === ClientStatus.PAUSED) {
        setState({ status: ClientStatus.PLAYING, message: '' });
        return;
    }

    if (state.status === ClientStatus.CONNECTING) {
        return;
    }

    if (state.status === ClientStatus.DISCONNECTED) {
        setState({ message: 'Server connection is closed. Refresh to reconnect.' });
        return;
    }

    const name = nameInput.value.trim();
    if (sendJoin(name)) {
        setState({ status: ClientStatus.CONNECTING, name, message: 'Joining...' });
    } else {
        setState({ status: ClientStatus.DISCONNECTED, message: 'Disconnected' });
    }
}

/**
 * @param {string} status
 */
function getPlayButtonLabel(status) {
    switch (status) {
        case ClientStatus.CONNECTING:
            return 'Joining...';
        case ClientStatus.PAUSED:
            return 'Resume';
        case ClientStatus.DEAD:
            return 'Play Again';
        case ClientStatus.DISCONNECTED:
            return 'Disconnected';
        default:
            return 'Play';
    }
}

/**
 * @param {EventTarget | null} target
 */
function isTypingTarget(target) {
    return target instanceof HTMLInputElement
        || target instanceof HTMLTextAreaElement
        || target instanceof HTMLSelectElement;
}

