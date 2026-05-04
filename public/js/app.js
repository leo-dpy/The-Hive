document.addEventListener('DOMContentLoaded', () => {
    const pingBtn = document.getElementById('pingBtn');
    const responseCard = document.getElementById('apiResponse');
    const jsonOutput = document.getElementById('jsonOutput');
    const statusIndicator = document.querySelector('.status-indicator');

    pingBtn.addEventListener('click', async () => {
        try {
            // Button loading state
            const originalText = pingBtn.textContent;
            pingBtn.textContent = 'Testing...';
            pingBtn.disabled = true;

            const response = await fetch('/api/ping');
            if (!response.ok) {
                throw new Error(`HTTP error: ${response.status}`);
            }
            const data = await response.json();
            
            // Display response
            responseCard.style.display = 'flex';
            statusIndicator.style.backgroundColor = 'var(--color-accent)';
            jsonOutput.textContent = `Server responded: ${data.message}`;

        } catch (error) {
            console.error("API Error:", error);
            responseCard.style.display = 'flex';
            statusIndicator.style.backgroundColor = '#ef4444'; // Red for error
            jsonOutput.textContent = "Connection to server failed.";
        } finally {
            // Restore button
            setTimeout(() => {
                pingBtn.textContent = 'Test Connectivity';
                pingBtn.disabled = false;
            }, 600);
        }
    });
});
