const contactsList = document.getElementById('contacts-list');
const chatPane = document.getElementById('chat-pane');
const chatEmpty = document.getElementById('chat-empty');
const chatHistory = document.getElementById('chat-history');
const chatInput = document.getElementById('chat-input');
const chatSend = document.getElementById('chat-send');
const chatAvatar = document.getElementById('chat-avatar');
const chatUsername = document.getElementById('chat-username');
const btnNewMsg = document.getElementById('btn-new-msg');
const modalNew = document.getElementById('modal-new');
const btnStartChat = document.getElementById('btn-start-chat');
const newMsgUsername = document.getElementById('new-msg-username');

let currentUserId = 0;
let activeChatId = 0;
let ws = null;

async function init() {
    const meRes = await fetch('/api/me');
    if (!meRes.ok) {
        window.location.href = '/';
        return;
    }
    const meData = await meRes.json();
    currentUserId = meData.id;
    document.getElementById('current-username').textContent = meData.username;

    loadConversations();
    connectWebSocket();
}

async function loadConversations() {
    try {
        const res = await fetch('/api/conversations');
        const convos = await res.json();
        
        contactsList.innerHTML = '';
        if (!convos || convos.length === 0) {
            contactsList.innerHTML = '<p style="text-align:center; color:#888; margin-top:20px; font-size:0.9rem;">Aucune conversation</p>';
            return;
        }

        convos.forEach(c => {
            const div = document.createElement('div');
            div.className = `contact-item ${c.user_id === activeChatId ? 'active' : ''}`;
            div.onclick = () => openChat(c.user_id, c.username, c.profile_picture);
            
            div.innerHTML = `
                <img src="${c.profile_picture}" class="contact-avatar">
                <div class="contact-info">
                    <div class="contact-name">${c.username}</div>
                    <div class="contact-preview">${escapeHTML(c.last_message || "Nouvelle conversation")}</div>
                </div>
            `;
            contactsList.appendChild(div);
        });
    } catch (e) {
        console.error(e);
    }
}

async function openChat(userId, username, avatar) {
    activeChatId = userId;
    chatEmpty.style.display = 'none';
    chatPane.style.display = 'flex';
    
    chatUsername.textContent = username;
    chatAvatar.src = avatar || '/uploads/avatars/default.png';
    chatHistory.innerHTML = '<p style="text-align:center; color:#888;">Chargement...</p>';

    // Mettre en surbrillance dans la liste
    loadConversations();

    try {
        const res = await fetch(`/api/messages?user_id=${userId}`);
        const messages = await res.json();
        chatHistory.innerHTML = '';
        
        if (messages && messages.length > 0) {
            messages.forEach(msg => appendMessage(msg));
        } else {
            chatHistory.innerHTML = '<p style="text-align:center; color:#888; margin-top:20px;">Début de la discussion</p>';
        }
        scrollToBottom();
    } catch (e) {
        console.error(e);
    }
}

function appendMessage(msg) {
    // Retirer le texte "Début de la discussion" s'il existe
    if (chatHistory.innerHTML.includes("Début de la discussion")) {
        chatHistory.innerHTML = '';
    }

    const div = document.createElement('div');
    const isSent = msg.sender_id === currentUserId;
    div.className = `bubble ${isSent ? 'bubble-sent' : 'bubble-received'}`;
    div.textContent = msg.content; // évite les injections XSS
    chatHistory.appendChild(div);
    scrollToBottom();
}

function scrollToBottom() {
    chatHistory.scrollTop = chatHistory.scrollHeight;
}

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${window.location.host}/api/ws`);
    
    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        if (msg.type === "new_dm") {
            if (msg.sender_id === activeChatId) {
                appendMessage(msg);
            }
            loadConversations(); // Met à jour l'aperçu du dernier message
        }
    };
    
    ws.onclose = () => {
        setTimeout(connectWebSocket, 3000); // Reconnexion automatique
    };
}

function sendMessage() {
    const text = chatInput.value.trim();
    if (!text || !activeChatId || !ws) return;

    const msg = {
        receiver_id: activeChatId,
        content: text
    };
    
    ws.send(JSON.stringify(msg));
    
    // Optimistic UI : ajouter la bulle directement
    appendMessage({
        sender_id: currentUserId,
        content: text
    });
    
    chatInput.value = '';
    loadConversations(); // Met à jour le snippet dans la liste
}

chatSend.addEventListener('click', sendMessage);
chatInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') sendMessage();
});

btnNewMsg.addEventListener('click', () => modalNew.style.display = 'flex');

btnStartChat.addEventListener('click', async () => {
    const username = newMsgUsername.value.trim().replace('@', '');
    if (!username) return;

    try {
        const res = await fetch(`/api/user?username=${username}`);
        if (!res.ok) {
            alert("Utilisateur introuvable");
            return;
        }
        const data = await res.json();
        modalNew.style.display = 'none';
        newMsgUsername.value = '';
        openChat(data.id, data.username, data.profile_picture);
    } catch (e) {
        alert("Erreur lors de la recherche");
    }
});

function escapeHTML(str) {
    const div = document.createElement('div');
    div.innerText = str;
    return div.innerHTML;
}

init();
