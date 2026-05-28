/** @type {Snapshot | null} snapshot */
let prev = null;

/**
 * @typedef {{id: string, owner_id: string, owner_name: string, x: number, y: number, radius: number, mass: number, color: string}} CellData
 * @typedef {{x: number, y: number, radius: number, mass: number, color: string}} PelletData
 * @typedef {{id: string, x: number, y: number, radius: number, mass: number, color: string}} EjectData
 * @typedef {{id: string, x: number, y: number, fed_count: number, radius: number, mass: number, color: string}} VirusData
 * @typedef {{cells: CellData[], pellets: PelletData[], ejects: EjectData[], viruses: VirusData[], me: string, score: number, tick: number}} Snapshot
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
    const ejects = s2.ejects.map(e => {
        const p = s1.ejects.find(pe => pe.id === e.id);
        if (!p) {
            return e;   // new eject, no prev
        }

        return {
            ...e,
            x: p.x + (e.x - p.x) * t,
            y: p.y + (e.y - p.y) * t
        };
    });
    const viruses = s2.viruses.map(v => {
        const p = s1.viruses.find(pv => pv.id === v.id);
        if (!p) {
            return v;   // new virus, no prev
        }

        return {
            ...v,
            x: p.x + (v.x - p.x) * t,
            y: p.y + (v.y - p.y) * t
        };
    });

    return {
        ...s1,
        cells: cells,
        ejects: ejects,
        viruses: viruses
    };
}
