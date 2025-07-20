# Real-time Chat Demo

This demo shows how WebSockets work with Go and HTMX using a simple chat application.

## Overview

The chat application demonstrates:
- Real-time bidirectional communication using WebSockets
- HTMX WebSocket extension for seamless HTML updates
- Go WebSocket server with gorilla/websocket
- Type-safe Templ templates
- Tailwind CSS styling

## Architecture

### Backend (Go)
- **WebSocket Handler**: `server/handlers/chat.go`
  - Upgrades HTTP connections to WebSocket
  - Manages client connections with a hub pattern
  - Broadcasts messages to all connected clients
  - HTML escaping for security

- **Chat Hub**: Central message broadcasting system
  - Client registration/unregistration
  - Message queuing and distribution
  - Concurrent-safe with mutexes

### Frontend (HTMX + Templ)
- **Chat Page**: `pages/chat.templ`
  - WebSocket connection via `ws-connect="/chat/ws"`
  - Form submission via `ws-send` attribute
  - Auto-scrolling message container
  - Input clearing after send

- **Real-time Updates**: Messages sent as HTML fragments
  - Uses HTMX out-of-band swaps (`hx-swap-oob="beforeend"`)
  - Automatic DOM updates without JavaScript

## Running the Demo

1. **Build the application**:
   ```bash
   make build
   ```

2. **Start the server**:
   ```bash
   ./bin/app
   ```

3. **Open multiple browser tabs**:
   - Navigate to `http://localhost:8082/chat`
   - Open 2-3 tabs to simulate multiple users

4. **Test real-time messaging**:
   - Type messages in one tab
   - See them appear instantly in all other tabs
   - Messages include username, timestamp, and content

## Development Mode

For development with live reload:

```bash
make dev
```

This runs:
- Air for Go live reload
- Templ template watching
- Tailwind CSS watching

## Key Features Demonstrated

### 1. WebSocket Connection
```html
<div hx-ext="ws" ws-connect="/chat/ws">
```

### 2. Message Sending
```html
<form ws-send hx-on::after-request="this.reset()">
```

### 3. Real-time Updates
```html
<div id="messages" hx-swap-oob="beforeend">
```

### 4. Auto-scroll
```javascript
document.body.addEventListener('htmx:oobAfterSwap', function(e) {
    if (e.detail.target.id === 'messages') {
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
    }
});
```

## Message Flow

1. **User types message** → Form with `ws-send`
2. **HTMX serializes form** → JSON sent to WebSocket
3. **Go handler receives** → Parses JSON message
4. **Server broadcasts** → HTML fragment to all clients
5. **HTMX receives HTML** → Out-of-band swap updates DOM
6. **Auto-scroll triggers** → Smooth scroll to bottom

## Security Features

- HTML escaping prevents XSS attacks
- Input validation (empty message filtering)
- CORS handling for WebSocket connections

## Extending the Demo

### Add User Authentication
- Modify `ChatMessage` struct to include user IDs
- Add session management
- Personalize usernames

### Add Chat Rooms
- Extend hub to support multiple rooms
- Add room selection UI
- Separate message broadcasting by room

### Add Message History
- Store messages in database
- Load recent messages on connection
- Add message persistence

### Add Rich Features
- File uploads
- Emoji support
- Message reactions
- Typing indicators

## Dependencies

- `github.com/gorilla/websocket` - WebSocket implementation
- `github.com/a-h/templ` - Type-safe templates
- HTMX WebSocket extension - Frontend WebSocket handling