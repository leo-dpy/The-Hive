const urlParams = new URLSearchParams(window.location.search);
const targetUsername = urlParams.get('username');

const headerUsername = document.getElementById('header-username');
const profileContainer = document.getElementById('profile-container');
const pAvatar = document.getElementById('p-avatar');
const pUsername = document.getElementById('p-username');
const pBio = document.getElementById('p-bio');
const pFollowers = document.getElementById('p-followers');
const pFollowing = document.getElementById('p-following');
const followBtn = document.getElementById('follow-btn');

const tabPosts = document.getElementById('tab-posts');
const tabLikes = document.getElementById('tab-likes');
const tabsContainer = document.getElementById('tabs-container');
const timeline = document.getElementById('timeline');
const loadingMsg = document.getElementById('loading-msg');

let currentProfileId = 0;
let currentTab = 'posts';

if (!targetUsername) {
    window.location.href = '/home.html';
}

async function loadProfile() {
    try {
        const res = await fetch(`/api/user?username=${targetUsername}`);
        if (!res.ok) {
            loadingMsg.textContent = "Utilisateur introuvable.";
            loadingMsg.style.color = "#e74c3c";
            return;
        }

        const data = await res.json();
        currentProfileId = data.id;

        // Mise à jour de l'UI
        headerUsername.textContent = `@${data.username}`;
        pUsername.textContent = `@${data.username}`;
        pBio.textContent = data.bio || "Cette abeille n'a pas encore de bio.";
        pAvatar.src = data.profile_picture;
        pFollowers.textContent = data.followers_count;
        pFollowing.textContent = data.following_count;
        
        profileContainer.style.display = 'block';
        tabsContainer.style.display = 'flex';

        // Check if viewing own profile
        const meRes = await fetch('/api/me');
        if (meRes.ok) {
            const meData = await meRes.json();
            if (meData.id !== data.id) {
                followBtn.style.display = 'block';
                updateFollowBtn(data.is_following);
                
                followBtn.onclick = async () => {
                    await fetch('/api/follow', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ target_id: data.id })
                    });
                    loadProfile(); // Recharger pour maj les compteurs
                };
            }
        }

        loadFeed();

    } catch (err) {
        console.error(err);
    }
}

function updateFollowBtn(isFollowing) {
    if (isFollowing) {
        followBtn.textContent = "Se désabonner";
        followBtn.style.backgroundColor = "transparent";
        followBtn.style.color = "var(--text-color)";
        followBtn.style.border = "1px solid var(--border-color)";
    } else {
        followBtn.textContent = "S'abonner";
        followBtn.style.backgroundColor = "var(--text-color)";
        followBtn.style.color = "white";
        followBtn.style.border = "none";
    }
}

async function loadFeed() {
    timeline.innerHTML = '<p style="text-align:center; color:#888; margin-top:40px;">Chargement...</p>';
    
    const endpoint = currentTab === 'posts' 
        ? `/api/user/posts?username=${targetUsername}` 
        : `/api/user/likes?username=${targetUsername}`;

    try {
        const res = await fetch(endpoint);
        const posts = await res.json();

        timeline.innerHTML = '';
        if (!posts || posts.length === 0) {
            timeline.innerHTML = `<p style="text-align:center; color:#888; margin-top:40px;">Aucun post à afficher.</p>`;
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

tabPosts.addEventListener('click', () => {
    currentTab = 'posts';
    tabPosts.classList.add('active');
    tabLikes.classList.remove('active');
    loadFeed();
});

tabLikes.addEventListener('click', () => {
    currentTab = 'likes';
    tabLikes.classList.add('active');
    tabPosts.classList.remove('active');
    loadFeed();
});

// Init
loadProfile();
