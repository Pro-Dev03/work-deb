const { contextBridge } = require('electron')

// exposing safe APIs to renderer process
contextBridge.exposeInMainWorld('electronAPI', {
  platform: process.platform,
  version: process.version,
  // يمكن إضافة المزيد من الـ APIs الآمنة هنا
})