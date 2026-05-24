/** @type {Snapshot | null} snapshot */
let prev = null;

/**
 * @typedef {{id: string, owner_id: string, x: number, y: number, radius: number, mass: number, color: string}} CellData
 * @typedef {{x: number, y: number, radius: number, mass: number, color: string}} PelletData
 * @typedef {{x: number, y: number, radius: number, mass: number, color: string}} VirusData
 * @typedef {{cells: CellData[], pellets: PelletData[], viruses: VirusData[], me: string, score: number, tick: number}} Snapshot
 */

/** @type {Snapshot | null} snapshot */
let curr = null;

let lastTickTime = 0;

const tickRate = 20;

/**
 * @param {Snapshot} snapshot
 */
export function pushSnapshot(snapshot) {
    prev = curr;
    curr = snapshot;
    lastTickTime = performance.now();
}

export function getInterpolated() {
    if (!prev || !curr) {
        return curr;
    }

    const t = Math.min(
        (performance.now() - lastTickTime) / (1000 / tickRate),
        1
    );
    return interpolateState(prev, curr, t);
}

/**
 * @param {Snapshot} s1
 * @param {Snapshot} s2
 * @param {number} t
 */
function interpolateState(s1, s2, t) {
    const cells = s2.cells.map(c => {
        const p = s1.cells.find(pc => pc.id === c.id);
        if (!p) {
            return c;   // new cell, no prev
        }

        return {
            ...c,
            x: p.x + (c.x - p.x) * t,
            y: p.y + (c.y - p.y) * t
        };
    });
    return { ...s1, cells: cells };
}
