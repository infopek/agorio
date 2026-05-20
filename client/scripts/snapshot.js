/**
 * @typedef {{owner_id: string, x: number, y: number, radius: number, mass: number, color: string}} CellData
 * @typedef {{x: number, y: number, radius: number, mass: number, color: string}} PelletData
 * @typedef {{x: number, y: number, radius: number, mass: number, color: string}} VirusData
 * @typedef {{cells: CellData[], pellets: PelletData[], viruses: VirusData[], me: string[], score: number, tick: number}} Snapshot
 */

/** @type {Snapshot | null} */
let current = null;

/**
 * @param {Snapshot} snapshot
 */
export function update(snapshot) {
    current = snapshot;
}

export function getSnapshot() {
    return current;
}
