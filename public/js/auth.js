let isLoginMode = true;

const form = document.getElementById('auth-form');
const toggleBtn = document.getElementById('toggle-mode');
const usernameGroup = document.getElementById('username-group');
const submitBtn = document.getElementById('submit-btn');
const errorMsg = document.getElementById('error-msg');

toggleBtn.addEventListener('click', () => {
    isLoginMode = !isLoginMode;
    if (isLoginMode) {
        usernameGroup.style.display = 'none';
        submitBtn.textContent = 'Se connecter';
        toggleBtn.textContent = "Pas encore de compte ? Créer un profil";
    } else {
        usernameGroup.style.display = 'block';
        submitBtn.textContent = "S'inscrire";
        toggleBtn.textContent = 'Déjà membre ? Se connecter';
    }
    errorMsg.textContent = '';
});

form.addEventListener('submit', async (e) => {
    e.preventDefault();
    errorMsg.textContent = '';

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const username = document.getElementById('username').value;

    const endpoint = isLoginMode ? '/api/login' : '/api/register';
    
    // Le backend de login accepte l'email dans le champ username
    const payload = isLoginMode 
        ? { username: email, password } 
        : { username, email, password };

    submitBtn.disabled = true;
    submitBtn.textContent = "Chargement...";

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

        // Redirection vers la Timeline
        window.location.href = '/home.html';
    } catch (err) {
        errorMsg.textContent = 'Impossible de joindre le serveur The Hive.';
        submitBtn.disabled = false;
        submitBtn.textContent = isLoginMode ? "Se connecter" : "S'inscrire";
    }
});
