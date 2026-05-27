document.addEventListener("DOMContentLoaded", async () => {
    // Attendre que checkAuth soit résolu (défini dans auth-modal.js)
    // On attend la vraie résolution de l'auth, pas un délai arbitraire
    if (typeof window.checkAuth === 'function') {
        try {
            await window.checkAuth();
        } catch (e) {
            console.warn('sidebar: checkAuth failed', e);
        }
    } else {
        // auth-modal.js n'est pas chargé sur cette page — tenter de fetch /api/me directement
        try {
            const res = await fetch('/api/me');
            if (res.ok) {
                window.currentUser = await res.json();
            }
        } catch (e) {
            console.warn('sidebar: fallback auth check failed', e);
        }
    }

    const profileLink = document.getElementById('my-profile-link');
    const messagesLink = document.querySelector('a[href="/messages.html"]');
    const settingsLink = document.querySelector('a[href="/settings.html"]');

    // Profil : mettre à jour le href si connecté
    if (profileLink) {
        if (window.currentUser) {
            profileLink.href = `/profile.html?username=${window.currentUser.username}`;
        } else {
            // Visiteur : cliquer sur "Profil" ouvre la modale de connexion
            profileLink.href = '#';
            profileLink.addEventListener('click', (e) => {
                if (!window.currentUser) {
                    e.preventDefault();
                    if (typeof window.requireAuth === 'function') {
                        requireAuth(() => {
                            window.location.href = `/profile.html?username=${window.currentUser.username}`;
                        });
                    } else {
                        // Pas de modale d'auth — rediriger vers home pour se connecter
                        window.location.href = '/home.html';
                    }
                }
            });
        }
    }

    // Messages : protégé par auth
    if (messagesLink) {
        messagesLink.addEventListener('click', (e) => {
            if (!window.currentUser) {
                e.preventDefault();
                if (typeof window.requireAuth === 'function') {
                    requireAuth(() => {
                        window.location.href = '/messages.html';
                    });
                } else {
                    window.location.href = '/home.html';
                }
            }
            // Si connecté, le lien href="/messages.html" fonctionne normalement
        });
    }

    // Paramètres : protégé par auth
    if (settingsLink) {
        settingsLink.addEventListener('click', (e) => {
            if (!window.currentUser) {
                e.preventDefault();
                if (typeof window.requireAuth === 'function') {
                    requireAuth(() => {
                        window.location.href = '/settings.html';
                    });
                } else {
                    window.location.href = '/home.html';
                }
            }
            // Si connecté, le lien href="/settings.html" fonctionne normalement
        });
    }
});
