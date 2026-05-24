import { Vec2 } from './vec2.js';
import { clamp } from './utils.js';

const MIN_ZOOM_OVERRIDE = 0.6;
const MAX_ZOOM_OVERRIDE = 1.3;

export class Camera {
    constructor() {
        this.position = new Vec2(0, 0);
        this.zoom = 1;
        this.zoomOverride = 1;  // scrolling
    }

    /**
     * @param {import('./interpolation.js').Snapshot} snapshot
     * @param {HTMLCanvasElement} canvas
     */
    update(snapshot, canvas) {
        let myCells = snapshot.cells.filter(c => c.owner_id == snapshot.me);
        if (myCells.length === 0) {
            return;
        }

        let sumX = 0.0
        let sumY = 0.0
        let sumMass = 0.0
        for (let i = 0; i < myCells.length; i++) {
            sumX += myCells[i].x * myCells[i].mass;
            sumY += myCells[i].y * myCells[i].mass;
            sumMass += myCells[i].mass;
        }

        const baseZoom = canvas.height / (200.0 + Math.sqrt(sumMass) * 33.0);
        //console.log('mass: ', sumMass);
        //console.log('baseZoom: ', baseZoom);

        this.position.x = sumX / sumMass;
        this.position.y = sumY / sumMass;
        this.zoom = baseZoom * this.zoomOverride;
    }

    /**
     * @param {number} scale
     */
    zoomBy(scale) {
        this.zoomOverride = clamp(this.zoomOverride * scale, MIN_ZOOM_OVERRIDE, MAX_ZOOM_OVERRIDE);
    }

    /**
     * @param {Vec2} position
     * @param {HTMLCanvasElement} canvas
     */
    worldToScreen(position, canvas) {
        return new Vec2(
            (position.x - this.position.x) * this.zoom + canvas.width / 2,
            (position.y - this.position.y) * this.zoom + canvas.height / 2
        );
    }

    /**
     * @param {number} radius
     */
    worldToScreenRadius(radius) {
        return radius * this.zoom;
    }
}

