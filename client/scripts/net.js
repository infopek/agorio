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
 * @param {string} name
 */
export function sendJoin(name) {
    if (ws === null) {
        return;
    }

    ws.send(JSON.stringify({
        t: 'join',
        name: name
    }));
}

/**
 * @param {number} x
 * @param {number} y
 */
export function sendMove(x, y) {
    if (ws === null) {
        return;
    }

    console.log('sending move with ', x, 'and ', y);
    ws.send(JSON.stringify({
        t: 'move',
        x: x,
        y: y
    }));
}

/**
 *
 */
export function sendSplit() {
    if (ws === null) {
        return;
    }

    console.log('sending split');
    ws.send(JSON.stringify({
        t: 'split',
    }));
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
