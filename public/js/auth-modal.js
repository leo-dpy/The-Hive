// ══════════════════════════════════════════════════════
//  AUTH MODAL — Browse-first, Login-on-action
//  Injects a login/register modal into any page.
//  Usage: requireAuth(() => { /* protected action */ });
// ══════════════════════════════════════════════════════

(function () {
    // ── Global auth state ──
    window.currentUser = null;

    // ── Check session on page load ──
    window.checkAuth = async function () {
        try {
            const res = await fetch('/api/me');
            if (res.ok) {
                window.currentUser = await res.json();
                document.body.classList.add('is-logged-in');
                document.body.classList.remove('is-guest');
                return window.currentUser;
            }
        } catch (e) {
            console.warn('Auth check failed:', e);
        }
        window.currentUser = null;
        document.body.classList.add('is-guest');
        document.body.classList.remove('is-logged-in');
        return null;
    };

    // ── Require auth gate ──
    // If logged in, executes callback immediately.
    // If not, shows the auth modal and executes callback after successful login.
    window.requireAuth = function (callback) {
        if (window.currentUser) {
            callback();
            return;
        }
        showAuthModal(callback);
    };

    // ── Pending callback after login ──
    let pendingCallback = null;

    // ── Inject modal HTML into DOM ──
    function createModal() {
        if (document.getElementById('auth-modal-overlay')) return;

        const overlay = document.createElement('div');
        overlay.id = 'auth-modal-overlay';
        overlay.className = 'auth-modal-overlay';
        overlay.innerHTML = `
            <div class="auth-modal" id="auth-modal">
                <button class="auth-modal-close" id="auth-modal-close" aria-label="Fermer">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                    </svg>
                </button>

                <div class="auth-modal-logo">The <span>Hive.</span></div>
                <p class="auth-modal-subtitle" id="auth-modal-subtitle">Connectez-vous pour interagir</p>
                <p class="auth-modal-error" id="auth-modal-error"></p>

                <form id="auth-modal-form">
                    <div class="form-group" id="auth-modal-username-group" style="display: none;">
                        <input type="text" id="auth-modal-username" placeholder="Pseudo (@nom)">
                    </div>
                    <div class="form-group">
                        <input type="email" id="auth-modal-email" placeholder="Adresse e-mail" required>
                    </div>
                    <div class="form-group">
                        <input type="password" id="auth-modal-password" placeholder="Mot de passe" required>
                    </div>
                    <button type="submit" class="btn-large auth-modal-submit" id="auth-modal-submit">Se connecter</button>
                </form>

                <div class="auth-modal-toggle" id="auth-modal-toggle">
                    Pas encore dans la Ruche ? Créer un profil
                </div>
            </div>
        `;
        document.body.appendChild(overlay);

        // ── Event listeners ──
        let isLoginMode = true;

        const closeBtn = document.getElementById('auth-modal-close');
        const toggleBtn = document.getElementById('auth-modal-toggle');
        const form = document.getElementById('auth-modal-form');
        const submitBtn = document.getElementById('auth-modal-submit');
        const errorMsg = document.getElementById('auth-modal-error');
        const usernameGroup = document.getElementById('auth-modal-username-group');
        const subtitleEl = document.getElementById('auth-modal-subtitle');

        // Close modal
        closeBtn.addEventListener('click', () => hideAuthModal());
        overlay.addEventListener('click', (e) => {
            if (e.target === overlay) hideAuthModal();
        });

        // Escape key
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape' && overlay.classList.contains('visible')) {
                hideAuthModal();
            }
        });

        // Toggle login/register
        toggleBtn.addEventListener('click', () => {
            isLoginMode = !isLoginMode;
            errorMsg.textContent = '';
            if (isLoginMode) {
                usernameGroup.style.display = 'none';
                submitBtn.textContent = 'Se connecter';
                toggleBtn.textContent = "Pas encore dans la Ruche ? Créer un profil";
                subtitleEl.textContent = "Connectez-vous pour interagir";
            } else {
                usernameGroup.style.display = 'block';
                submitBtn.textContent = "S'inscrire";
                toggleBtn.textContent = 'Déjà membre ? Se connecter';
                subtitleEl.textContent = "Rejoignez la Ruche";
            }
        });

        // Submit form
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            errorMsg.textContent = '';

            const email = document.getElementById('auth-modal-email').value;
            const password = document.getElementById('auth-modal-password').value;
            const username = document.getElementById('auth-modal-username').value;

            const endpoint = isLoginMode ? '/api/login' : '/api/register';
            const payload = isLoginMode
                ? { username: email, password }
                : { username, email, password };

            submitBtn.disabled = true;
            submitBtn.textContent = 'Chargement...';

            try {
                const res = await fetch(endpoint, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(payload)
                });

                const data = await res.json();

                if (!res.ok) {
                    errorMsg.textContent = data.error || 'Une erreur est survenue';
                    submitBtn.disabled = false;
                    submitBtn.textContent = isLoginMode ? "Se connecter" : "S'inscrire";
                    return;
                }

                // Login successful — refresh user state
                await window.checkAuth();
                hideAuthModal();

                // Execute the pending action
                if (pendingCallback) {
                    const cb = pendingCallback;
                    pendingCallback = null;
                    cb();
                }

                // Refresh page-specific content if available
                if (typeof window.onAuthSuccess === 'function') {
                    window.onAuthSuccess();
                }

            } catch (err) {
                errorMsg.textContent = 'Impossible de joindre le serveur The Hive.';
                submitBtn.disabled = false;
                submitBtn.textContent = isLoginMode ? "Se connecter" : "S'inscrire";
            }
        });
    }

    function showAuthModal(callback) {
        createModal();
        pendingCallback = callback || null;

        const overlay = document.getElementById('auth-modal-overlay');
        const errorMsg = document.getElementById('auth-modal-error');
        const submitBtn = document.getElementById('auth-modal-submit');

        // Reset state
        if (errorMsg) errorMsg.textContent = '';
        if (submitBtn) {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Se connecter';
        }

        // Show with animation
        requestAnimationFrame(() => {
            overlay.classList.add('visible');
        });
    }

    function hideAuthModal() {
        const overlay = document.getElementById('auth-modal-overlay');
        if (overlay) {
            overlay.classList.remove('visible');
        }
        pendingCallback = null;
    }

    // Expose for direct usage
    window.showAuthModal = showAuthModal;
    window.hideAuthModal = hideAuthModal;

    // Auto-check auth on load
    document.addEventListener('DOMContentLoaded', () => {
        window.checkAuth();
    });
})();
