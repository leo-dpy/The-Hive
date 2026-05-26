document.addEventListener("DOMContentLoaded", async () => {
    // Attendre que checkAuth soit résolu (défini dans auth-modal.js)
    // Small delay to ensure auth-modal.js has run
    await new Promise(resolve => setTimeout(resolve, 50));
    
    const profileLink = document.getElementById('my-profile-link');
    const messagesLink = document.querySelector('a[href="/messages.html"]');
    const settingsLink = document.querySelector('a[href="/settings.html"]');

    if (profileLink) {
        if (window.currentUser) {
            profileLink.href = `/profile.html?username=${window.currentUser.username}`;
        } else {
            // Visiteur : cliquer sur "Profil" ouvre la modale de connexion
            profileLink.href = '#';
            profileLink.addEventListener('click', (e) => {
                if (!window.currentUser) {
                    e.preventDefault();
                    requireAuth(() => {
                        window.location.href = `/profile.html?username=${window.currentUser.username}`;
                    });
                }
            });
        }
    }

    // Messages : protégé par auth
    if (messagesLink) {
        messagesLink.addEventListener('click', (e) => {
            if (!window.currentUser) {
                e.preventDefault();
                requireAuth(() => {
                    window.location.href = '/messages.html';
                });
            }
        });
    }

    // Paramètres : protégé par auth
    if (settingsLink) {
        settingsLink.addEventListener('click', (e) => {
            if (!window.currentUser) {
                e.preventDefault();
                requireAuth(() => {
                    window.location.href = '/settings.html';
                });
            }
        });
    }
});
