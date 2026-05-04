const form = document.getElementById('settings-form');
const usernameInput = document.getElementById('username');
const emailInput = document.getElementById('email');
const bioInput = document.getElementById('bio');
const newPasswordInput = document.getElementById('new_password');
const msg = document.getElementById('msg');
const logoutBtn = document.getElementById('logout-btn');
const avatarUpload = document.getElementById('avatar-upload');
const currentAvatar = document.getElementById('current-avatar');

async function loadProfile() {
    try {
        const res = await fetch('/api/me');
        if (!res.ok) {
            window.location.href = '/';
            return;
        }
        const data = await res.json();
        usernameInput.value = data.username || '';
        emailInput.value = data.email || '';
        bioInput.value = data.bio || '';
        currentAvatar.src = data.profile_picture || '/uploads/avatars/default.png';
    } catch (err) {
        msg.textContent = "Erreur de connexion";
        msg.style.color = "#e74c3c";
    }
}

form.addEventListener('submit', async (e) => {
    e.preventDefault();
    msg.textContent = "Sauvegarde en cours...";
    msg.style.color = "#888";

    const payload = {
        username: usernameInput.value.trim(),
        email: emailInput.value.trim(),
        bio: bioInput.value.trim(),
        new_password: newPasswordInput.value
    };

    try {
        const res = await fetch('/api/user/update', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        const data = await res.json();

        if (!res.ok) {
            msg.textContent = data.error || "Erreur lors de la sauvegarde";
            msg.style.color = "#e74c3c";
        } else {
            msg.textContent = data.message || "Profil mis à jour !";
            msg.style.color = "#2ecc71";
            newPasswordInput.value = ''; // On vide le mot de passe après sauvegarde
        }
    } catch (err) {
        msg.textContent = "Erreur de communication avec le serveur";
        msg.style.color = "#e74c3c";
    }
});

logoutBtn.addEventListener('click', async () => {
    await fetch('/api/logout', { method: 'POST' });
    window.location.href = '/';
});

avatarUpload.addEventListener('change', async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    msg.textContent = "Upload de l'avatar en cours...";
    msg.style.color = "#888";

    const formData = new FormData();
    formData.append('avatar', file);

    try {
        const res = await fetch('/api/user/avatar', {
            method: 'POST',
            body: formData
        });
        const data = await res.json();
        if (res.ok) {
            currentAvatar.src = data.url;
            msg.textContent = "Avatar mis à jour !";
            msg.style.color = "#2ecc71";
        } else {
            msg.textContent = data.error || "Erreur d'upload";
            msg.style.color = "#e74c3c";
        }
    } catch (err) {
        msg.textContent = "Erreur réseau";
        msg.style.color = "#e74c3c";
    }
});

// Chargement initial
loadProfile();
