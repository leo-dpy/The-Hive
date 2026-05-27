const timeline = document.getElementById('timeline');
const postContent = document.getElementById('post-content');
const publishBtn = document.getElementById('publish-btn');
const categorySelect = document.getElementById('category-select');

// Charger les catégories pour le sélecteur
async function loadCategories() {
    try {
        const res = await fetch('/api/categories');
        const categories = await res.json();
        
        if (categorySelect && categories.length > 0) {
            categories.forEach(cat => {
                const option = document.createElement('option');
                option.value = cat.id;
                option.textContent = cat.name;
                categorySelect.appendChild(option);
            });
        }
    } catch (e) {
        console.error('Erreur chargement catégories:', e);
    }
}

async function loadFeed() {
    try {
        let posts = [];
        // Si connecté, charger le feed "following" d'abord
        if (window.currentUser) {
            const followingRes = await fetch('/api/feed/following');
            if (followingRes.ok) {
                posts = await followingRes.json();
            }
            // Si le feed "following" est vide (nouveau compte, personne suivi),
            // on charge le feed public pour que la page ne soit jamais vide
            if (!posts || posts.length === 0) {
                const publicRes = await fetch('/api/posts');
                if (publicRes.ok) {
                    posts = await publicRes.json();
                }
            }
        } else {
            const res = await fetch('/api/posts');
            if (res.ok) {
                posts = await res.json();
            }
        }
        
        timeline.innerHTML = '';
        if (!posts || posts.length === 0) {
            timeline.innerHTML = `
                <div style="text-align:center; margin-top: 40px;">
                    <p style="color:var(--text-muted); font-size:1.1rem; margin-bottom: 15px;">
                        Aucun post pour le moment. Soyez le premier à publier !
                    </p>
                    <button class="btn-small" onclick="document.getElementById('post-content').focus()">Écrire un post</button>
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
                            <span class="action" onclick="reactWithAuth(${post.id}, 1)">▲ ${post.likes}</span>
                            <span class="action" onclick="reactWithAuth(${post.id}, -1)">▼ ${post.dislikes}</span>
                            <span class="action" onclick="window.location.href='/explore.html?post_id=${post.id}'"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="vertical-align:-2px;margin-right:4px;"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>${post.comment_count}</span>
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

publishBtn.addEventListener('click', () => {
    requireAuth(async () => {
        const content = postContent.value.trim();
        if (!content) return;

        publishBtn.disabled = true;

        // Récupérer la catégorie sélectionnée
        const categoryId = categorySelect ? parseInt(categorySelect.value) : null;

        const body = { content };
        if (categoryId) {
            body.category_id = categoryId;
        }

        try {
            const res = await fetch('/api/posts', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body)
            });
            
            if (res.ok) {
                postContent.value = '';
                if (categorySelect) categorySelect.value = '';
                loadFeed(); // Recharger le feed pour voir le nouveau post
            }
        } catch (err) {
            console.error(err);
        }
        publishBtn.disabled = false;
    });
});

// React with auth check
window.reactWithAuth = (postId, value) => {
    requireAuth(async () => {
        await fetch('/api/react', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ target_type: 'post', target_id: postId, value })
        });
        loadFeed();
    });
};

// Callback appelé après connexion réussie via la modale
window.onAuthSuccess = () => {
    loadFeed();
};

// Chargement initial — attendre que checkAuth soit terminé
document.addEventListener('DOMContentLoaded', async () => {
    await window.checkAuth();
    loadCategories();
    loadFeed();
});
