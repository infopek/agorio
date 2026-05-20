export class Vec2 {
    /**
     * @param {number} x
     * @param {number} y
     */
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }

    /**
     * @param {Vec2} other
     */
    add(other) {
        return new Vec2(
            this.x + other.x,
            this.y + other.y
        );
    }

    /**
     * @param {Vec2} other
     */
    sub(other) {
        return new Vec2(
            this.x - other.x,
            this.y - other.y
        );
    }

    /**
     * @param {number} s
     */
    scale(s) {
        return new Vec2(
            this.x * s,
            this.y * s
        );
    }

    normalize() {
        let mag = this.magnitude();
        if (mag == 0.0) {
            return new Vec2(0, 0);
        }

        return new Vec2(
            this.x / mag,
            this.y / mag
        );
    }

    magnitude() {
        return Math.sqrt(this.x * this.x + this.y * this.y);
    }

    magnitudeSq() {
        return this.x * this.x + this.y * this.y;
    }

    /**
     * @param {Vec2} other
     */
    distanceTo(other) {
        let temp = other.sub(this);
        return temp.magnitude();
    }

    /**
     * @param {Vec2}   other
     * @param {number} t
     */
    lerp(other, t) {
        return new Vec2(
            this.x * t + other.x * (1.0 - t),
            this.y * t + other.y * (1.0 - t),
        )
    }
}
