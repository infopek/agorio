export const ClientStatus = Object.freeze({
    MENU: 'menu',
    CONNECTING: 'connecting',
    PAUSED: 'paused',
    PLAYING: 'playing',
    DEAD: 'dead',
    DISCONNECTED: 'disconnected',
});

/** @type {{status: string, name: string, message: string}} */
let state = {
    status: ClientStatus.MENU,
    name: '',
    message: '',
};

/** @type {Set<(s: state) => void>} */
const listeners = new Set();

export function getState() {
    return state;
}

/** setState
 *
 * @param {Object} patch
 *
 */
export function setState(patch) {
    state = {
        ...state,
        ...patch,
    };

    for (const listener of listeners) {
        listener(state);
    }
}

/** subscribe
 *
 * @param {(s: state) => void} listener
 *
 */
export function subscribe(listener) {
    listeners.add(listener);
    listener(state);

    return () => {
        listeners.delete(listener);
    }
}


