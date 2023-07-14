const input = document.querySelector('#textarea')
const messagesEl = document.querySelector('#messages')
const username = document.querySelector('#username')
const send = document.querySelector('#send')

const url = "ws://" + window.location.host + "/ws"
const ws = new WebSocket(url)
const chatID = window.location.pathname.match(/\/chats\/(\d+)/)[1]

ws.onmessage = (msg) => insertMessage(JSON.parse(msg.data))

send.onclick = () => {
    const message = {
        username: username.value,
        content: input.value,
        chatID: chatID,
    }

    console.log(message)

    ws.send(JSON.stringify(message))
    input.value = ""
}

function insertMessage(message) {
    const newMessageDiv = document.createElement("div")
    newMessageDiv.setAttribute('class', 'chat-message')

    newMessageDiv.textContent = `${message.username}: ${message.content}`

    messagesEl.prepend(newMessageDiv)
}