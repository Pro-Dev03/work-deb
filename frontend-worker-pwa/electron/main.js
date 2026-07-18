import { app, BrowserWindow, Menu, session, ipcMain } from 'electron'
import path from 'path'
import { fileURLToPath } from 'url'
import WebSocket from 'ws'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const isDev = process.env.NODE_ENV === 'development' || !app.isPackaged

// تعطيل sandbox لتجنب مشاكل الصلاحيات
app.commandLine.appendSwitch('--no-sandbox')
app.commandLine.appendSwitch('--disable-setuid-sandbox')
// إعدادات إضافية للسماح باتصالات WebSocket عبر HTTPS
app.commandLine.appendSwitch('--ignore-certificate-errors')
app.commandLine.appendSwitch('--allow-running-insecure-content')
app.commandLine.appendSwitch('--disable-features=VizDisplayCompositor')
// إعدادات إضافية لتحسين الاستقرار على Windows
app.commandLine.appendSwitch('--disable-software-rasterizer')
app.commandLine.appendSwitch('--in-process-gpu')
// إعدادات إضافية لتحسين التوافق مع Linux Wayland
if (process.platform === 'linux') {
  app.commandLine.appendSwitch('--ozone-platform=x11')
}

let mainWindow
let wsClient = null

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    minWidth: 360,
    minHeight: 600,
    show: false, // إخفاء النافذة حتى تحميل المحتوى
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      enableRemoteModule: false,
      webSecurity: false, // إيقاف webSecurity للسماح بالاتصال بالـ API الخارجي
      preload: path.join(__dirname, 'preload.js'),
      additionalArguments: ['--disable-web-security', '--allow-running-insecure-content']
    },
    icon: path.join(__dirname, '../dist/icon.png'), // استخدام المسار الصحيح للأيقونة
    title: 'WorkTrack - تطبيق الموظف'
  })

  // إظهار النافذة عندما تكون جاهزة
  mainWindow.once('ready-to-show', () => {
    mainWindow.show()
  })

  // تحميل التطبيق
  if (isDev) {
    mainWindow.loadURL('http://localhost:3000')
    mainWindow.webContents.openDevTools()
  } else {
    mainWindow.loadFile(path.join(__dirname, '../dist/index.html'))
    // تمرير متغير البيئة للكشف عن Electron
    mainWindow.webContents.executeJavaScript('window.isElectron = true')
  }

  // معالجة أخطاء التحميل
  mainWindow.webContents.on('did-fail-load', (event, errorCode, errorDescription, validatedURL) => {
    console.error('Failed to load:', errorCode, errorDescription, validatedURL)
  })

  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

// إنشاء القائمة
function createMenu() {
  const template = [
    {
      label: 'ملف',
      submenu: [
        {
          label: 'تحديث الصفحة',
          accelerator: 'CmdOrCtrl+R',
          click: () => {
            mainWindow.webContents.reload()
          }
        },
        {
          label: 'خروج',
          accelerator: 'CmdOrCtrl+Q',
          click: () => {
            app.quit()
          }
        }
      ]
    },
    {
      label: 'عرض',
      submenu: [
        {
          label: 'تكبير',
          accelerator: 'CmdOrCtrl+Plus',
          click: () => {
            mainWindow.webContents.setZoomLevel(mainWindow.webContents.getZoomLevel() + 0.5)
          }
        },
        {
          label: 'تصغير',
          accelerator: 'CmdOrCtrl+-',
          click: () => {
            mainWindow.webContents.setZoomLevel(mainWindow.webContents.getZoomLevel() - 0.5)
          }
        },
        {
          label: 'إعادة التعيين',
          accelerator: 'CmdOrCtrl+0',
          click: () => {
            mainWindow.webContents.setZoomLevel(0)
          }
        },
        { type: 'separator' },
        {
          label: 'مطور',
          accelerator: 'F12',
          click: () => {
            mainWindow.webContents.toggleDevTools()
          }
        }
      ]
    },
    {
      label: 'مساعدة',
      submenu: [
        {
          label: 'حول WorkTrack',
          click: () => {
            // يمكن إضافة نافذة عن التطبيق هنا
          }
        }
      ]
    }
  ]

  const menu = Menu.buildFromTemplate(template)
  Menu.setApplicationMenu(menu)
}

// WebSocket handlers
ipcMain.handle('websocket-connect', (event, url) => {
  try {
    if (wsClient) {
      wsClient.close()
    }

    wsClient = new WebSocket(url)

    wsClient.on('open', () => {
      console.log('✅ WebSocket connected in main process')
      if (mainWindow) {
        mainWindow.webContents.send('websocket-open')
      }
    })

    wsClient.on('message', (data) => {
      try {
        const message = JSON.parse(data.toString())
        if (mainWindow) {
          mainWindow.webContents.send('websocket-message', message)
        }
      } catch (e) {
        console.error('Error parsing WebSocket message:', e)
      }
    })

    wsClient.on('error', (error) => {
      console.error('WebSocket error:', error)
      if (mainWindow) {
        mainWindow.webContents.send('websocket-error', error.message)
      }
    })

    wsClient.on('close', (code, reason) => {
      console.log('WebSocket closed:', code, reason.toString())
      if (mainWindow) {
        mainWindow.webContents.send('websocket-close', code, reason.toString())
      }
    })

    return { success: true }
  } catch (error) {
    console.error('Failed to connect WebSocket:', error)
    return { success: false, error: error.message }
  }
})

ipcMain.handle('websocket-send', (event, data) => {
  if (wsClient && wsClient.readyState === WebSocket.OPEN) {
    wsClient.send(JSON.stringify(data))
    return { success: true }
  }
  return { success: false, error: 'WebSocket not connected' }
})

ipcMain.handle('websocket-disconnect', () => {
  if (wsClient) {
    wsClient.close()
    wsClient = null
  }
  return { success: true }
})

app.whenReady().then(() => {
  // إعداد CORS وحل مشاكل الشبكة
  session.defaultSession.webRequest.onHeadersReceived((details, callback) => {
    callback({
      responseHeaders: {
        ...details.responseHeaders,
        'Access-Control-Allow-Origin': ['*'],
        'Access-Control-Allow-Methods': ['GET, POST, PUT, DELETE, OPTIONS'],
        'Access-Control-Allow-Headers': ['Content-Type, Authorization']
      }
    })
  })
  
  // إعداد SSL/CORS للسماح بالاتصال بـ API خارجي
  session.defaultSession.webRequest.onBeforeSendHeaders((details, callback) => {
    callback({
      requestHeaders: {
        ...details.requestHeaders,
        'User-Agent': 'WorkTrack-Electron-App/1.0'
      }
    })
  })
  
  // إعداد WebSocket للسماح بالاتصالات عبر بروتوكول ws:// و wss://
  session.defaultSession.webRequest.onBeforeRequest((details, callback) => {
    const url = details.url
    // السماح بجميع الاتصالات بما فيها WebSocket
    callback({ cancel: false })
  })
  
  // إعداد إضافي لضمان عمل WebSocket مع HTTPS
  session.defaultSession.setCertificateVerifyProc((request, callback) => {
    // تجاهل أخطاء الشهادات للسماح بالاتصال بالخوادم الخارجية
    callback(0)
  })
  
  createWindow()
  createMenu()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow()
    }
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// التعامل مع الأخطاء
process.on('uncaughtException', (error) => {
  console.error('Uncaught Exception:', error)
})

app.on('renderer-process-crashed', (event, webContents, killed) => {
  console.error('Renderer process crashed:', killed)
  app.relaunch()
  app.exit(1)
})

// إضافة تسجيل إضافي لتصحيح الأخطاء
app.on('web-contents-created', (event, contents) => {
  contents.on('console-message', (event, level, message, line, sourceId) => {
    console.log(`[Console ${level}] ${message}`)
  })
})