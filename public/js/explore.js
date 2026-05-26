const mainContent = document.getElementById('main-content');
const breadcrumb = document.getElementById('breadcrumb');
const sidebarCategories = document.getElementById('sidebar-categories');
const widgetCategories = document.getElementById('widget-categories');

// ── État de la navigation ──
let allCategories = [];

// ── Initialisation ──
async function init() {
    // Charger les catégories pour la sidebar
    try {
        const res = await fetch('/api/categories');
        allCategories = await res.json();
        renderSidebarCategories();
    } catch (e) {
        console.error('Erreur chargement catégories:', e);
    }

    // Vérifier les paramètres URL
    const params = new URLSearchParams(window.location.search);
    const postId = params.get('post_id');
    const categoryId = params.get('category_id');

    if (postId) {
        showPost(parseInt(postId));
    } else if (categoryId) {
        const cat = allCategories.find(c => c.id === parseInt(categoryId));
        showCategory(parseInt(categoryId), cat ? cat.name : 'Catégorie');
    } else {
        showCategories();
    }
}

// ── Sidebar catégories ──
function renderSidebarCategories() {
    if (allCategories.length === 0) return;
    widgetCategories.style.display = 'block';
    sidebarCategories.innerHTML = '';
    allCategories.forEach(cat => {
        const div = document.createElement('div');
        div.className = 'sidebar-cat-item';
        div.textContent = cat.name;
        div.onclick = () => navigateToCategory(cat.id, cat.name);
        sidebarCategories.appendChild(div);
    });
}

// ── Breadcrumb ──
function setBreadcrumb(items) {
    breadcrumb.innerHTML = '';
    items.forEach((item, i) => {
        const span = document.createElement('span');
        span.className = 'breadcrumb-item' + (i === items.length - 1 ? ' active' : '');
        span.textContent = item.label;
        if (item.action && i < items.length - 1) {
            span.style.cursor = 'pointer';
            span.onclick = item.action;
        }
        breadcrumb.appendChild(span);

        if (i < items.length - 1) {
            const sep = document.createElement('span');
            sep.className = 'breadcrumb-sep';
            sep.textContent = '›';
            breadcrumb.appendChild(sep);
        }
    });
}

// ── Navigation helpers ──
function navigateToCategory(id, name) {
    history.pushState(null, '', `/explore.html?category_id=${id}`);
    showCategory(id, name);
}

function navigateToPost(id) {
    history.pushState(null, '', `/explore.html?post_id=${id}`);
    showPost(id);
}

function navigateToHome() {
    history.pushState(null, '', '/explore.html');
    showCategories();
}

// ── VUE 1 : Liste des catégories ──
function showCategories() {
    setBreadcrumb([{ label: 'Explorer' }]);

    mainContent.innerHTML = '';
    
    if (allCategories.length === 0) {
        mainContent.innerHTML = '<p style="text-align:center; color:#888; margin-top:40px;">Aucune catégorie trouvée.</p>';
        return;
    }

    const grid = document.createElement('div');
    grid.className = 'categories-grid';

    allCategories.forEach(cat => {
        const card = document.createElement('div');
        card.className = 'category-card';
        card.onclick = () => navigateToCategory(cat.id, cat.name);
        
        // Icône par catégorie
        const icons = {
            'Général': '💬',
            'Technologie': '💻',
            'Programmation': '💻',
            'Gaming': '🎮',
            'Jeux Vidéo': '🎮',
            'Musique': '🎵',
            'Sport': '⚽',
            'Cinéma & Séries': '🎬',
            'Études & Carrière': '🎓',
            'Art & Design': '🎨',
            'Cuisine': '🍳',
            'Voyage': '🌍',
            'Science': '🔬',
            'Actualités': '📰'
        };
        const icon = icons[cat.name] || '📁';

        card.innerHTML = `
            <div class="category-icon">${icon}</div>
            <div class="category-card-name">${cat.name}</div>
            <div class="category-card-desc">${cat.description || ''}</div>
        `;
        grid.appendChild(card);
    });

    mainContent.appendChild(grid);
}

// ── VUE 2 : Posts d'une catégorie ──
async function showCategory(categoryId, categoryName) {
    setBreadcrumb([
        { label: 'Explorer', action: navigateToHome },
        { label: categoryName }
    ]);

    mainContent.innerHTML = '<p style="text-align:center; color:#888; margin-top:40px;">Chargement des sujets...</p>';

    try {
        const res = await fetch(`/api/posts?category_id=${categoryId}`);
        const posts = await res.json();

        mainContent.innerHTML = '';

        if (!posts || posts.length === 0) {
            mainContent.innerHTML = `
                <div style="text-align:center; margin-top:40px;">
                    <p style="color:var(--text-muted); font-size:1.1rem; margin-bottom:15px;">Aucun sujet dans cette catégorie.</p>
                    <button class="btn-small" onclick="window.location.href='/home.html'">Créer un sujet</button>
                </div>`;
            return;
        }

        posts.forEach(post => {
            const date = new Date(post.created_at).toLocaleString('fr-FR', {
                day: 'numeric', month: 'short', hour: '2-digit', minute:'2-digit'
            });
            const div = document.createElement('div');
            div.className = 'post thread-item';
            div.onclick = (e) => {
                // Ne pas naviguer si on clique sur un lien ou bouton
                if (e.target.closest('a') || e.target.closest('.action')) return;
                navigateToPost(post.id);
            };
            div.style.cursor = 'pointer';
            div.innerHTML = `
                <div style="display:flex; gap:15px; align-items:flex-start;">
                    <a href="/profile.html?username=${post.author}">
                        <img src="${post.profile_picture}" style="width:45px; height:45px; border-radius:50%; object-fit:cover; border: 1px solid var(--border-color);">
                    </a>
                    <div style="flex:1;">
                        <div class="post-author">
                            <a href="/profile.html?username=${post.author}" style="color:inherit; text-decoration:none;">@${post.author}</a>
                            <span style="font-size:12px; font-weight:normal; color:#888; margin-left:8px;">${date}</span>
                        </div>
                        ${post.title ? `<div class="post-title">${escapeHTML(post.title)}</div>` : ''}
                        <div class="post-content">${escapeHTML(post.content)}</div>
                        <div class="post-actions">
                            <span class="action" onclick="event.stopPropagation(); react(${post.id}, 1)">▲ ${post.likes}</span>
                            <span class="action" onclick="event.stopPropagation(); react(${post.id}, -1)">▼ ${post.dislikes}</span>
                            <span class="action" onclick="event.stopPropagation(); navigateToPost(${post.id})">💬 ${post.comment_count}</span>
                        </div>
                    </div>
                </div>
            `;
            mainContent.appendChild(div);
        });
    } catch (err) {
        console.error('Erreur chargement posts:', err);
        mainContent.innerHTML = '<p style="text-align:center; color:#e74c3c;">Erreur de chargement.</p>';
    }
}

// ── VUE 3 : Détail d'un post + commentaires ──
async function showPost(postId) {
    mainContent.innerHTML = '<p style="text-align:center; color:#888; margin-top:40px;">Chargement de la discussion...</p>';

    try {
        // Charger le post
        const postRes = await fetch(`/api/post?id=${postId}`);
        if (!postRes.ok) {
            mainContent.innerHTML = '<p style="text-align:center; color:#e74c3c; margin-top:40px;">Post introuvable.</p>';
            return;
        }
        const post = await postRes.json();

        // Trouver la catégorie pour le breadcrumb
        // On doit récupérer le category_id depuis l'API - pour l'instant on utilise "Discussion"
        setBreadcrumb([
            { label: 'Explorer', action: navigateToHome },
            { label: post.title || 'Discussion' }
        ]);

        // Charger les commentaires
        const commentsRes = await fetch(`/api/comments?post_id=${postId}`);
        const comments = await commentsRes.json();

        const date = new Date(post.created_at).toLocaleString('fr-FR', {
            day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute:'2-digit'
        });

        mainContent.innerHTML = `
            <div class="post-detail">
                <div class="post-detail-header">
                    <a href="/profile.html?username=${post.author}">
                        <img src="${post.profile_picture}" class="post-detail-avatar">
                    </a>
                    <div>
                        <a href="/profile.html?username=${post.author}" class="post-detail-author">@${post.author}</a>
                        <div class="post-detail-date">${date}</div>
                    </div>
                </div>
                ${post.title ? `<h2 class="post-detail-title">${escapeHTML(post.title)}</h2>` : ''}
                <div class="post-detail-content">${escapeHTML(post.content)}</div>
                <div class="post-detail-actions">
                    <span class="action" onclick="react(${post.id}, 1)">▲ ${post.likes}</span>
                    <span class="action" onclick="react(${post.id}, -1)">▼ ${post.dislikes}</span>
                    <span class="action">💬 ${post.comment_count}</span>
                </div>
            </div>

            <div class="comments-section">
                <div class="comments-header">
                    <h3>Commentaires (${comments.length})</h3>
                </div>
                
                <div class="comment-form" id="comment-form">
                    <textarea id="comment-input" placeholder="Ajouter un commentaire..." rows="3" onfocus="handleCommentFocus(event)"></textarea>
                    <div style="display:flex; justify-content:flex-end; margin-top:10px;">
                        <button class="btn-small" id="comment-submit">Commenter</button>
                    </div>
                </div>

                <div id="comments-list">
                    ${comments.length === 0 
                        ? '<p class="no-comments">Aucun commentaire. Soyez le premier à réagir !</p>'
                        : comments.map(c => renderComment(c)).join('')
                    }
                </div>
            </div>
        `;

        // Gérer la soumission du commentaire
        document.getElementById('comment-submit').addEventListener('click', () => submitComment(postId));
        document.getElementById('comment-input').addEventListener('keydown', (e) => {
            if (e.key === 'Enter' && e.ctrlKey) submitComment(postId);
        });

    } catch (err) {
        console.error('Erreur chargement post:', err);
        mainContent.innerHTML = '<p style="text-align:center; color:#e74c3c;">Erreur de chargement.</p>';
    }
}

function renderComment(c) {
    const date = new Date(c.created_at).toLocaleString('fr-FR', {
        day: 'numeric', month: 'short', hour: '2-digit', minute:'2-digit'
    });
    return `
        <div class="comment">
            <a href="/profile.html?username=${c.author}">
                <img src="${c.profile_picture}" class="comment-avatar">
            </a>
            <div class="comment-body">
                <div class="comment-meta">
                    <a href="/profile.html?username=${c.author}" class="comment-author">@${c.author}</a>
                    <span class="comment-date">${date}</span>
                </div>
                <div class="comment-content">${escapeHTML(c.content)}</div>
            </div>
        </div>
    `;
}

async function submitComment(postId) {
    requireAuth(async () => {
        const input = document.getElementById('comment-input');
        const content = input.value.trim();
        if (!content) return;

        const btn = document.getElementById('comment-submit');
        btn.disabled = true;
        btn.textContent = 'Envoi...';

        try {
            const res = await fetch('/api/comments', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ post_id: postId, content })
            });

            if (res.ok) {
                input.value = '';
                // Recharger la vue du post pour voir le nouveau commentaire
                showPost(postId);
            } else {
                const data = await res.json();
                alert(data.error || "Erreur lors de l'envoi du commentaire.");
            }
        } catch (e) {
            alert("Erreur réseau");
        }
        btn.disabled = false;
        btn.textContent = 'Commenter';
    });
}

// ── Utilitaires ──
function escapeHTML(str) {
    const div = document.createElement('div');
    div.innerText = str;
    return div.innerHTML;
}

window.react = (postId, value) => {
    requireAuth(async () => {
        await fetch('/api/react', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ target_type: 'post', target_id: postId, value })
        });
        // Recharger la vue actuelle
        const params = new URLSearchParams(window.location.search);
        if (params.get('post_id')) {
            showPost(parseInt(params.get('post_id')));
        } else if (params.get('category_id')) {
            const catId = parseInt(params.get('category_id'));
            const cat = allCategories.find(c => c.id === catId);
            showCategory(catId, cat ? cat.name : 'Catégorie');
        }
    });
};

// Handle comment textarea focus — require auth
window.handleCommentFocus = function(e) {
    if (!window.currentUser) {
        e.target.blur();
        requireAuth(() => {
            const input = document.getElementById('comment-input');
            if (input) input.focus();
        });
    }
};

// Gestion du bouton retour du navigateur
window.addEventListener('popstate', () => {
    const params = new URLSearchParams(window.location.search);
    const postId = params.get('post_id');
    const categoryId = params.get('category_id');

    if (postId) {
        showPost(parseInt(postId));
    } else if (categoryId) {
        const cat = allCategories.find(c => c.id === parseInt(categoryId));
        showCategory(parseInt(categoryId), cat ? cat.name : 'Catégorie');
    } else {
        showCategories();
    }
});

// Lancement
init();
