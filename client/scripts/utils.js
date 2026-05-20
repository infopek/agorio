/**
* @param {number} num
* @param {number} min
* @param {number} max
*/
export function clamp(num, min, max) {
    return Math.min(Math.max(num, min), max);
}
