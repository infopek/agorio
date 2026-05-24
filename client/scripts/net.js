import { pushSnapshot } from './interpolation.js';
import { toggleMenu } from './input.js';

/** @type {WebSocket | null} */
let ws = null;
let alive = false;

export function connect() {
    ws = new WebSocket("ws://localhost:8080/ws");
    ws.onopen = () => console.log("connected");
    ws.onmessage = (msg) => {
        const data = JSON.parse(msg.data);
        if (data.t === 'snapshot') {
            pushSnapshot(data);
        } else if (data.t === 'death') {
            alive = false;
            toggleMenu(true);   // show menu
        }
    }
    ws.onerror = (e) => console.error('error: ', e);
    ws.onclose = (e) => console.error('closed: ', e.code, e.reason);
}

/** send
 *
 * @param {Object} obj
 */
function send(obj) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(obj));
    }
}

/**
 *
 * @param {string} name
 */
export function sendJoin(name) {
    send({
        t: 'join',
        name: name
    });
    alive = true;
}

/** sendMove
 *
 * @param {number} x
 * @param {number} y
 */
export function sendMove(x, y) {
    send({
        t: 'move',
        x: x,
        y: y
    });
}

/** sendSplit
 *
 */
export function sendSplit() {
    send({
        t: 'split'
    });
}

/** sendFeed
 *
 */
export function sendFeed() {
    send({
        t: 'feed',
    });
}
