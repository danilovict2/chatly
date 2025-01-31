const pusher = new Pusher('398ed8ab4241f6a50dec', {
    cluster: 'eu',
    channelAuthorization: {
        endpoint: '/pusher/auth',
        headers: { 'X-CSRF-Token': JSON.parse(document.getElementById('CSRF').textContent) },
    },
});

const setSelectedUserStatus = () => {
    const statusComponent = document.getElementById('user-status');
    if (statusComponent) {
        const userID = JSON.parse(statusComponent.getAttribute('user-id'));
        const isUserOnline = isOnline(userID);
        statusComponent.innerText = isUserOnline ? "Online" : "Offline";
        
        const onlineCircle = document.getElementById("online-circle");
        if (onlineCircle) {
            onlineCircle.classList.toggle('hidden', !isUserOnline);
        }
    }
}

const initializeRealTimeMessageHandler = (authUserUsername, selectedUserUsername) => {
    const channel = pusher.subscribe('message');
    // Receive messages sent to authenticated user
    channel.bind(`to.${authUserUsername}`, message => {
        if (message.Sender !== selectedUserUsername) {
            return;
        }

        const messageHTML = document.createElement('div');
        messageHTML.classList.add("chat", "chat-start");
        messageHTML.innerHTML = `
            <div class="chat-image avatar">
				<div class="size-10 rounded-full border">
					<img
						src=${message.SenderAvatar}
						alt=${message.Sender}
					/>
				</div>
			</div>
			<div class="chat-header mb-1">
				<time class="text-xs opacity-50 ml-1">
					${message.CreatedAt}
				</time>
			</div>
			<div class="chat-bubble flex flex-col">
				${message.Image !== "" ?
                `<img
                        src="/public/img/${message.Image}"
                        alt="Attachment"
                        className="sm:max-w-[200px] rounded-md mb-2"
                    />`
                : ``}
				<p>${message.Text}</p>
			</div>
		</div>
        `;

        document.getElementById('chat-container').appendChild(messageHTML);
        scrollToBottom();
    });
}

const isOnline = (userID) => {
    const user = onlineUsersChannel.members.get(userID)
    return user !== null
} 

const onlineUsersChannel = pusher.subscribe('presence-users');
onlineUsersChannel.bind('pusher:subscription_succeeded', setSelectedUserStatus);
onlineUsersChannel.bind('pusher:member_added', setSelectedUserStatus);
onlineUsersChannel.bind('pusher:member_removed', setSelectedUserStatus);