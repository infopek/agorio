/**
 * @param {{id: string, name: string, score: number}[]} entries
 * @param {string} myId
 */
export function updateLeaderboard(entries, myId) {
    const el = document.getElementById('leaderboard');
    if (!el) {
        return;
    }

    let html = '<h3>Leaderboard</h3><ol>';
    for (const entry of entries) {
        const cls = entry.id === myId ? ' class="me"' : ''; // highlight ourselves
        html += `<li${cls}>${entry.name}: ${entry.score}</li>`;
    }
    html += '</ol>';
    el.innerHTML = html;
}
