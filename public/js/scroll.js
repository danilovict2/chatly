const scrollToBottom = () => {
    console.log('here')
    const container = document.getElementById('chat');
    container.scrollTop = container.scrollHeight;
};

document.getElementById("chat").addEventListener("load", scrollToBottom())