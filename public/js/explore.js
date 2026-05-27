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
        
        // Icône SVG par catégorie
        const icons = {
            'Général': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>',
            'Technologie': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>',
            'Programmation': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 18 22 12 16 6"></polyline><polyline points="8 6 2 12 8 18"></polyline></svg>',
            'Gaming': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><line x1="6" y1="12" x2="10" y2="12"></line><line x1="8" y1="10" x2="8" y2="14"></line><line x1="15" y1="13" x2="15.01" y2="13"></line><line x1="18" y1="11" x2="18.01" y2="11"></line><path d="M17.32 5H6.68a4 4 0 0 0-3.978 3.59c-.006.052-.01.101-.017.152C2.604 9.416 2 14.456 2 16a3 3 0 0 0 3 3c1 0 1.5-.5 2-1l1.414-1.414A2 2 0 0 1 9.828 16h4.344a2 2 0 0 1 1.414.586L17 18c.5.5 1 1 2 1a3 3 0 0 0 3-3c0-1.544-.604-6.584-.685-7.258-.007-.05-.011-.1-.017-.151A4 4 0 0 0 17.32 5z"></path></svg>',
            'Jeux Vidéo': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><line x1="6" y1="12" x2="10" y2="12"></line><line x1="8" y1="10" x2="8" y2="14"></line><line x1="15" y1="13" x2="15.01" y2="13"></line><line x1="18" y1="11" x2="18.01" y2="11"></line><path d="M17.32 5H6.68a4 4 0 0 0-3.978 3.59c-.006.052-.01.101-.017.152C2.604 9.416 2 14.456 2 16a3 3 0 0 0 3 3c1 0 1.5-.5 2-1l1.414-1.414A2 2 0 0 1 9.828 16h4.344a2 2 0 0 1 1.414.586L17 18c.5.5 1 1 2 1a3 3 0 0 0 3-3c0-1.544-.604-6.584-.685-7.258-.007-.05-.011-.1-.017-.151A4 4 0 0 0 17.32 5z"></path></svg>',
            'Musique': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M9 18V5l12-2v13"></path><circle cx="6" cy="18" r="3"></circle><circle cx="18" cy="16" r="3"></circle></svg>',
            'Sport': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M12 2a14.5 14.5 0 0 0 0 20 14.5 14.5 0 0 0 0-20"></path><path d="M2 12h20"></path></svg>',
            'Cinéma & Séries': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="20" rx="2.18" ry="2.18"></rect><line x1="7" y1="2" x2="7" y2="22"></line><line x1="17" y1="2" x2="17" y2="22"></line><line x1="2" y1="12" x2="22" y2="12"></line><line x1="2" y1="7" x2="7" y2="7"></line><line x1="2" y1="17" x2="7" y2="17"></line><line x1="17" y1="17" x2="22" y2="17"></line><line x1="17" y1="7" x2="22" y2="7"></line></svg>',
            'Études & Carrière': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M22 10v6M2 10l10-5 10 5-10 5z"></path><path d="M6 12v5c0 2 2 3 6 3s6-1 6-3v-5"></path></svg>',
            'Art & Design': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="13.5" cy="6.5" r=".5" fill="currentColor"></circle><circle cx="17.5" cy="10.5" r=".5" fill="currentColor"></circle><circle cx="8.5" cy="7.5" r=".5" fill="currentColor"></circle><circle cx="6.5" cy="12.5" r=".5" fill="currentColor"></circle><path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10c.926 0 1.648-.746 1.648-1.688 0-.437-.18-.835-.437-1.125-.29-.289-.438-.652-.438-1.125a1.64 1.64 0 0 1 1.668-1.668h1.996c3.051 0 5.555-2.503 5.555-5.554C21.965 6.012 17.461 2 12 2z"></path></svg>',
            'Cuisine': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8h1a4 4 0 0 1 0 8h-1"></path><path d="M2 8h16v9a4 4 0 0 1-4 4H6a4 4 0 0 1-4-4V8z"></path><line x1="6" y1="1" x2="6" y2="4"></line><line x1="10" y1="1" x2="10" y2="4"></line><line x1="14" y1="1" x2="14" y2="4"></line></svg>',
            'Voyage': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>',
            'Science': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M6 4h8a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"></path><path d="M6 12h9a4 4 0 0 1 4 4 4 4 0 0 1-4 4H6z"></path></svg>',
            'Actualités': '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M19 20H5a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v1"></path><path d="M21 12a9 9 0 0 0-9-9"></path><path d="M21 12H3"></path><path d="M14 2v4h4"></path><rect x="14" y="14" width="8" height="8" rx="1"></rect></svg>'
        };
        const icon = icons[cat.name] || '<svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>';

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
                            <span class="action" onclick="event.stopPropagation(); navigateToPost(${post.id})"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="vertical-align:-2px;margin-right:4px;"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>${post.comment_count}</span>
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
                    <span class="action"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="vertical-align:-2px;margin-right:4px;"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>${post.comment_count}</span>
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
