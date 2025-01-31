const scrollToBottom = () => {
    const container = document.getElementById('chat-container');
    if (container) {
        // Wait for images to load
        setTimeout(() => {
            container.scrollTop = container.scrollHeight;
        }, 450);
    }
};