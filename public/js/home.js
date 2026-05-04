const timeline = document.getElementById('timeline');
const postContent = document.getElementById('post-content');
const publishBtn = document.getElementById('publish-btn');

// Initialiser le WebSocket pour être prêt à recevoir des DMs (optionnel pour l'instant)
// const ws = new WebSocket(`ws://${window.location.host}/api/ws`);
// ws.onmessage = (e) => { console.log("Nouveau DM reçu:", JSON.parse(e.data)); };

async function loadFeed() {
    try {
        const res = await fetch('/api/feed/following');
        if (!res.ok) {
            window.location.href = '/';
            return;
        }
        const posts = await res.json();
        
        timeline.innerHTML = '';
        if (!posts || posts.length === 0) {
            timeline.innerHTML = `
                <div style="text-align:center; margin-top: 40px;">
                    <p style="color:var(--text-muted); font-size:1.1rem; margin-bottom: 15px;">Vous ne suivez personne, ou vos abonnements n'ont rien publié.</p>
                    <button class="btn-small" onclick="window.location.href='/explore.html'">Découvrir des Abeilles</button>
                </div>
            `;
            return;
        }

        posts.forEach(post => {
            const date = new Date(post.created_at).toLocaleString('fr-FR', {
                day: 'numeric', month: 'short', hour: '2-digit', minute:'2-digit'
            });
            const div = document.createElement('div');
            div.className = 'post';
            div.innerHTML = `
                <div style="display:flex; gap:15px; align-items:flex-start;">
                    <a href="/profile.html?username=${post.author}">
                        <img src="${post.profile_picture}" style="width:45px; height:45px; border-radius:50%; object-fit:cover; border: 1px solid var(--border-color);">
                    </a>
                    <div style="flex:1;">
                        <div class="post-author"><a href="/profile.html?username=${post.author}" style="color:inherit; text-decoration:none;">@${post.author}</a> <span style="font-size:12px; font-weight:normal; color:#888; margin-left:8px;">${date}</span></div>
                        <div class="post-content">${escapeHTML(post.content)}</div>
                        <div class="post-actions">
                            <span class="action" onclick="react(${post.id}, 1)">▲ ${post.likes}</span>
                            <span class="action" onclick="react(${post.id}, -1)">▼ ${post.dislikes}</span>
                            <span class="action">💬 ${post.comment_count}</span>
                        </div>
                    </div>
                </div>
            `;
            timeline.appendChild(div);
        });
    } catch (err) {
        console.error(err);
        timeline.innerHTML = '<p style="text-align:center; color:#e74c3c;">Impossible de charger le feed.</p>';
    }
}

// Fonction utilitaire pour empêcher le XSS basique
function escapeHTML(str) {
    const div = document.createElement('div');
    div.innerText = str;
    return div.innerHTML;
}

publishBtn.addEventListener('click', async () => {
    const content = postContent.value.trim();
    if (!content) return;

    publishBtn.disabled = true;
    try {
        const res = await fetch('/api/posts', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ content })
        });
        
        if (res.ok) {
            postContent.value = '';
            loadFeed(); // Recharger le feed pour voir le nouveau post
        }
    } catch (err) {
        console.error(err);
    }
    publishBtn.disabled = false;
});

// Exposer la fonction globale pour les clics sur Upvote/Downvote
window.react = async (postId, value) => {
    await fetch('/api/react', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ target_type: 'post', target_id: postId, value })
    });
    loadFeed();
};

// Chargement initial
loadFeed();
