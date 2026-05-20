import { Vec2 } from './vec2.js';
import { clamp } from './utils.js';

const MIN_ZOOM = 0.5;
const MAX_ZOOM = 2.0;

export class Camera {
    constructor() {
        this.position = new Vec2(0, 0);
        this.zoom = 1;
    }

    /**
     * @param {import('./snapshot.js').Snapshot} snapshot
     * @param {HTMLCanvasElement} canvas
     */
    update(snapshot, canvas) {
        let sumX = 0.0
        let sumY = 0.0
        let sumMass = 0.0
        for (let i = 0; i < snapshot.cells.length; i++) {
            sumX += snapshot.cells[i].x * snapshot.cells[i].mass;
            sumY += snapshot.cells[i].y * snapshot.cells[i].mass;
            sumMass += snapshot.cells[i].mass;
        }

        if (sumMass === 0.0) {
            return;
        }

        this.position.x = sumX / sumMass;
        this.position.y = sumY / sumMass;
        this.zoom = clamp(canvas.height / (120.0 + sumMass * 0.3), MIN_ZOOM, MAX_ZOOM);
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

