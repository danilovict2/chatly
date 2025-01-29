const enableMessageSubscription = () => {
    const authUserUsername = JSON.parse(document.getElementById('senderUsername').textContent);
    const selectedUserUsername = JSON.parse(document.getElementById('receiverUsername').textContent);

    const pusher = new Pusher('398ed8ab4241f6a50dec', {
        cluster: 'eu'
    });

    const channel = pusher.subscribe('messages');
    // Receive messages sent to authenticated user
    channel.bind(`to-${authUserUsername}`, message => {
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

        document.getElementById("chat-container").appendChild(messageHTML);
        scrollToBottom();
    });
}

enableMessageSubscription();