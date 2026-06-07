import { pushSnapshot } from './interpolation.js';
import { updateLeaderboard, clearLeaderboard } from './leaderboard.js';
import { ClientStatus, getState, setState } from './state.js';
import { Config } from './config.js';

/** @type {WebSocket | null} */
let ws = null;

export function connect() {
    ws = new WebSocket(Config.websocketUrl);

    ws.onopen = () => {
        clearLeaderboard();
        setState({ status: ClientStatus.MENU, message: '' });
    };

    ws.onclose = () => {
        clearLeaderboard();
        setState({ status: ClientStatus.DISCONNECTED, message: 'Disconnected' });
    };

    ws.onmessage = (msg) => {
        const data = JSON.parse(msg.data);
        if (data.t === 'snapshot') {
            pushSnapshot(data);

            if (getState().status === ClientStatus.CONNECTING) {
                setState({ status: ClientStatus.PLAYING, message: '' });
            }
        } else if (data.t === 'death') {
            setState({ status: ClientStatus.DEAD, message: 'You died' });
        } else if (data.t === 'leaderboard') {
            if (getState().status !== ClientStatus.MENU
                && getState().status !== ClientStatus.DISCONNECTED) {
                updateLeaderboard(data.entries, data.me);
            }
        }
    }
    ws.onerror = (e) => console.error('error: ', e);
}

/** send
 *
 * @param {Object} obj
 */
function send(obj) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify(obj));
        return true;
    }
    return false;
}

/**
 *
 * @param {string} name
 */
export function sendJoin(name) {
    return send({
        t: 'join',
        name: name
    });
}

/** sendMove
 *
 * @param {number} x
 * @param {number} y
 */
export function sendMove(x, y) {
    return send({
        t: 'move',
        x: x,
        y: y
    });
}

/** sendSplit
 *
 */
export function sendSplit() {
    return send({
        t: 'split'
    });
}

/** sendFeed
 *
 */
export function sendFeed() {
    return send({
        t: 'feed',
    });
}
