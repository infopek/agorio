import { update } from './snapshot.js';

/** @type {WebSocket | null} */
let ws = null;

export function connect() {
    ws = new WebSocket("ws://localhost:8080/ws");
    ws.onopen = () => console.log("connected");
    ws.onmessage = (msg) => {
        const data = JSON.parse(msg.data);
        if (data.t === 'snapshot') {
            update(data);
        }
    }
    ws.onerror = (e) => console.error('error: ', e);
    ws.onclose = (e) => console.error('closed: ', e.code, e.reason);
}

/**
 * @param {Object} obj
 */
function send(obj) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(obj));
    }
}

/**
 * @param {string} name
 */
export function sendJoin(name) {
    send({
        t: 'join',
        name: name
    });
}

/**
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

/**
 *
 */
export function sendSplit() {
    console.log('sending split');
    send({
        t: 'split'
    });
}

/**
 *
 */
export function sendFeed() {
    if (ws === null) {
        return;
    }

    console.log('sending feed');
    ws.send(JSON.stringify({
        t: 'feed',
    }));
}
