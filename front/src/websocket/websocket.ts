import type { WSMessage } from "./messages";

type MessageHandler<T extends WSMessage> = (data: T) => void;

const env = import.meta.env.VITE_APP_ENV;
const API_HOST = env == "production" ? document.location.origin : "http://localhost:8888";

class WebSocketService {
    private socket: WebSocket | null = null;
    private handlers: Map<string, MessageHandler<WSMessage>[]> = new Map();

    onConnect(callback?: () => void) {
        if (this.socket) return;
        const url = `${API_HOST}/ws`;
        
        this.socket = new WebSocket(url);
        console.log(`Connecting to WebSocket at ${url}`);

        this.socket.onopen = () => {
            console.log("WebSocket connection established");
            callback?.();
        };

        this.socket.onclose = () => {
            console.log("WebSocket connection closed");
            this.socket = null;
        };

        this.socket.onmessage = (event) => {
            const msg = JSON.parse(event.data);
            const handlers = this.handlers.get(msg.type);
            if (handlers) {
                handlers.forEach(h => h(msg));
            }
        };

        this.socket.onerror = (error) => {
            console.error("WebSocket error:", error);
            this.socket = null;
        };
    }

    sendMessage(message: WSMessage) {
        if (!this.socket) return;
        this.socket.send(JSON.stringify(message));
    }

    subscribe(topic: string) {
        const message: WSMessage = {
            type: 'subscribe',
            data: topic,
        };
        this.sendMessage(message);
    }

    unsubscribe(topic: string) {
        const message: WSMessage = {
            type: 'unsubscribe',
            data: topic,
        };
        this.sendMessage(message);
    }

    onMessage<T extends WSMessage>(type: string, handler: MessageHandler<T>) {
        const handlers = this.handlers.get(type) || [];
        handlers.push(handler as MessageHandler<WSMessage>);
        this.handlers.set(type, handlers);
    }

    offMessage<T extends WSMessage>(type: string, handler: MessageHandler<T>) {
        const handlers = this.handlers.get(type);
        if (handlers) {
            this.handlers.set(type, handlers.filter(h => h !== handler));
        }
    }

    clearHandlers() {
        this.handlers.clear();
    }

    disconnect() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }
}

export const wsService = new WebSocketService();
