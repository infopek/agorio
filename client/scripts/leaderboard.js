/**
 * @param {{id: string, name: string, score: number}[]} entries
 * @param {string} myId
 */
export function updateLeaderboard(entries, myId) {
    const el = document.getElementById('leaderboard');
    if (!el) {
        return;
    }

    const title = document.createElement('h3');
    title.textContent = 'Leaderboard';

    const list = document.createElement('ol');

    for (const entry of entries) {
        const item = document.createElement('li');

        const name = document.createElement('span');
        name.className = 'leaderboard-name';
        name.textContent = entry.name;

        const score = document.createElement('span');
        score.className = 'leaderboard-score';
        score.textContent = String(entry.score);

        item.append(name, score);

        if (entry.id === myId) {
            item.classList.add('me');
        }

        list.append(item);
    }

    el.replaceChildren(title, list);
}

export function clearLeaderboard() {
    const el = document.getElementById('leaderboard');
    if (!el) {
        return;
    }

    el.replaceChildren();
}
