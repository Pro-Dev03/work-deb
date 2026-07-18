const { app, BrowserWindow, Menu, session } = require('electron')
const path = require('path')
const isDev = require('electron-is-dev')

let mainWindow

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    minWidth: 1024,
    minHeight: 768,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      enableRemoteModule: false,
      webSecurity: false, // إيقاف webSecurity للسماح بالاتصال بالـ API الخارجي
      preload: path.join(__dirname, 'preload.js'),
      additionalArguments: ['--disable-web-security', '--allow-running-insecure-content']
    },
    icon: path.join(__dirname, '../public/favicon.ico'),
    title: 'WorkTrack - بوابة العملاء'
  })

  // تحميل التطبيق
  if (isDev) {
    mainWindow.loadURL('http://localhost:3002')
    mainWindow.webContents.openDevTools()
  } else {
    mainWindow.loadFile(path.join(__dirname, '../dist/index.html'))
    // تمرير متغير البيئة للكشف عن Electron
    mainWindow.webContents.executeJavaScript('window.isElectron = true')
  }

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