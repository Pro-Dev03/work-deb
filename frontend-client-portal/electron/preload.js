const { contextBridge, ipcRenderer } = require('electron')

// exposing safe APIs to renderer process
contextBridge.exposeInMainWorld('electronAPI', {
  platform: process.platform,
  version: process.version,
  isElectron: true,
  // WebSocket API للاتصال عبر Electron
  websocket: {
    connect: (url) => ipcRenderer.invoke('websocket-connect', url),
    send: (data) => ipcRenderer.invoke('websocket-send', data),
    disconnect: () => ipcRenderer.invoke('websocket-disconnect'),
    onMessage: (callback) => ipcRenderer.on('websocket-message', (event, data) => callback(data)),
    onOpen: (callback) => ipcRenderer.on('websocket-open', () => callback()),
    onError: (callback) => ipcRenderer.on('websocket-error', (event, error) => callback(error)),
    onClose: (callback) => ipcRenderer.on('websocket-close', (event, code, reason) => callback(code, reason))
  }
})