document.addEventListener("DOMContentLoaded", async () => {
    const profileLink = document.getElementById('my-profile-link');
    if (profileLink) {
        try {
            const res = await fetch('/api/me');
            if (res.ok) {
                const data = await res.json();
                profileLink.href = `/profile.html?username=${data.username}`;
            }
        } catch (err) {
            console.error(err);
        }
    }
});
