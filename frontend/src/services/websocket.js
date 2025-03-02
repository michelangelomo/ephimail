export default class WebSocketService {
    constructor() {
      this.socket = null;
      this.connected = false;
      this.reconnectAttempts = 0;
      this.maxReconnectAttempts = 5;
      this.reconnectTimeout = null;
      this.eventListeners = {
        'new_email': [],
        'email_deleted': [],
        'connect': [],
        'disconnect': [],
        'error': []
      };
    }
  
    /**
     * Connect to the WebSocket server
     */
    connect() {
      if (this.socket) {
        return;
      }
  
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${process.env.VUE_APP_BACKEND_URL.replace(/^https?:\/\//, '')}/ws`;
  
      this.socket = new WebSocket(wsUrl);
  
      this.socket.onopen = () => {
        console.log('WebSocket connected');
        this.connected = true;
        this.reconnectAttempts = 0;
        this._triggerEvent('connect');
      };
  
      this.socket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          if (message.type && this.eventListeners[message.type]) {
            this._triggerEvent(message.type, message.payload);
          }
        } catch (error) {
          console.error('WebSocket message error:', error);
        }
      };
  
      this.socket.onclose = () => {
        console.log('WebSocket disconnected');
        this.connected = false;
        this.socket = null;
        this._triggerEvent('disconnect');
  
        // Attempt to reconnect
        this._reconnect();
      };
  
      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        this._triggerEvent('error', error);
      };
    }
  
    /**
     * Disconnect from the WebSocket server
     */
    disconnect() {
      if (this.socket) {
        this.socket.close();
        this.socket = null;
        this.connected = false;
      }
  
      // Clear any pending reconnect
      if (this.reconnectTimeout) {
        clearTimeout(this.reconnectTimeout);
        this.reconnectTimeout = null;
      }
    }
  
    /**
     * Subscribe to an email inbox
     * @param {string} email - The email address to subscribe to
     */
    subscribeToInbox(email) {
      if (!this.connected || !this.socket) {
        // Queue the subscription for when we connect
        this.addEventListener('connect', () => {
          this.subscribeToInbox(email);
        }, { once: true });
        
        // Try to connect if not already
        if (!this.socket) {
          this.connect();
        }
        return;
      }
  
      const message = {
        type: 'subscribe',
        payload: {
          email: email
        }
      };
  
      this.socket.send(JSON.stringify(message));
    }
  
    /**
     * Unsubscribe from an email inbox
     * @param {string} email - The email address to unsubscribe from
     */
    unsubscribeFromInbox(email) {
      if (!this.connected || !this.socket) {
        return;
      }
  
      const message = {
        type: 'unsubscribe',
        payload: {
          email: email
        }
      };
  
      this.socket.send(JSON.stringify(message));
    }
  
    /**
     * Add an event listener
     * @param {string} event - The event name
     * @param {function} callback - The callback function
     * @param {object} options - Options (e.g., { once: true })
     */
    addEventListener(event, callback, options = {}) {
      if (!this.eventListeners[event]) {
        this.eventListeners[event] = [];
      }
  
      this.eventListeners[event].push({
        callback,
        options
      });
    }
  
    /**
     * Remove an event listener
     * @param {string} event - The event name
     * @param {function} callback - The callback function to remove
     */
    removeEventListener(event, callback) {
      if (!this.eventListeners[event]) {
        return;
      }
  
      this.eventListeners[event] = this.eventListeners[event].filter(
        listener => listener.callback !== callback
      );
    }
  
    /**
     * Trigger an event
     * @private
     */
    _triggerEvent(event, payload) {
      if (!this.eventListeners[event]) {
        return;
      }
  
      const listeners = [...this.eventListeners[event]];
      const onceListeners = [];
  
      listeners.forEach(listener => {
        listener.callback(payload);
        
        if (listener.options && listener.options.once) {
          onceListeners.push(listener);
        }
      });
  
      // Remove once listeners
      onceListeners.forEach(listener => {
        this.removeEventListener(event, listener.callback);
      });
    }
  
    /**
     * Attempt to reconnect to the WebSocket server
     * @private
     */
    _reconnect() {
      if (this.reconnectAttempts >= this.maxReconnectAttempts) {
        console.log('Max reconnect attempts reached');
        return;
      }
  
      this.reconnectAttempts++;
      const backoff = Math.min(30, Math.pow(2, this.reconnectAttempts)) * 1000;
      
      console.log(`Reconnecting in ${backoff / 1000} seconds...`);
      
      this.reconnectTimeout = setTimeout(() => {
        console.log('Attempting to reconnect...');
        this.connect();
      }, backoff);
    }
  }
  
  // Create a singleton instance
  const webSocketService = new WebSocketService();
  export { webSocketService };