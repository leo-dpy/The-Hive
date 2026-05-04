const timeline = document.getElementById('timeline');

async function loadFeed() {
    try {
        const res = await fetch('/api/posts');
        if (!res.ok) {
            window.location.href = '/';
            return;
        }
        const posts = await res.json();
        
        timeline.innerHTML = '';

        if (!posts || posts.length === 0) {
            timeline.innerHTML = '<p style="text-align:center; color:#888; margin-top: 20px;">Aucun post sur le réseau.</p>';
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
        console.error("Erreur chargement feed", err);
    }
}

function escapeHTML(str) {
    const div = document.createElement('div');
    div.innerText = str;
    return div.innerHTML;
}

window.react = async (postId, value) => {
    await fetch('/api/react', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ target_type: 'post', target_id: postId, value })
    });
    loadFeed();
};

loadFeed();
