const scrollToBottom = () => {
    const container = document.getElementById('chat-container');
    container.scrollTop = container.scrollHeight;
};

document.getElementById('chat-container').addEventListener('load', scrollToBottom())