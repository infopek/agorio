export const Config = Object.freeze({
    // Net
    websocketUrl: 'ws://localhost:8080/ws',

    // Interpolation
    tickRate: 20,

    // World / Rendering
    worldWidth: 15000,
    worldHeight: 15000,
    gridSize: 30,
    virusFeedToShoot: 7,

    // Camera
    minZoomOverride: 0.8,
    maxZoomOverride: 1.2,
    wheelZoomOutScale: 0.9,
    wheelZoomInScale: 1.1,
    cameraBaseViewSize: 200,
    cameraMassZoomFactor: 20,
    cameraPositionSmoothing: 0.1,
    cameraZoomSmoothing: 0.05,
});
