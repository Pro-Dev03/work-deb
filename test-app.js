const puppeteer = require('puppeteer');
const http = require('http');

// دالة للتحقق من أن الخادم يعمل
function checkServer(url) {
  return new Promise((resolve) => {
    const req = http.get(url, (res) => {
      resolve(true);
    });
    req.on('error', () => resolve(false));
    req.setTimeout(3000, () => {
      req.destroy();
      resolve(false);
    });
  });
}

async function testApp() {
  console.log('🚀 بدء اختبار التطبيق...\n');

  // التحقق من الخوادم
  console.log('📡 التحقق من الخوادم...');
  const frontendOk = await checkServer('http://localhost:3001');
  const backendOk = await checkServer('http://localhost:8080/health');
  
  console.log(`   Frontend (3001): ${frontendOk ? '✅ يعمل' : '❌ لا يعمل'}`);
  console.log(`   Backend (8080): ${backendOk ? '✅ يعمل' : '❌ لا يعمل'}`);
  
  if (!frontendOk || !backendOk) {
    console.log('\n❌ الخوادم لا تعمل، يرجى تشغيلها أولاً');
    return;
  }

  console.log('\n🌐 فتح المتصفح...');
  const browser = await puppeteer.launch({
    headless: false,
    defaultViewport: null,
    args: ['--start-maximized', '--no-sandbox', '--disable-setuid-sandbox']
  });

  try {
    const page = await browser.newPage();
    
    // تفعيل تسجيل أخطاء الكونسول
    page.on('console', msg => {
      if (msg.type() === 'error') {
        console.log('❌ خطأ في الكونسول:', msg.text());
      }
    });

    page.on('pageerror', error => {
      console.log('❌ خطأ في الصفحة:', error.message);
    });

    // 1. فتح صفحة تسجيل الدخول
    console.log('\n📄 فتح صفحة تسجيل الدخول...');
    await page.goto('http://localhost:3001/login', { waitUntil: 'networkidle2' });
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    // فحص عناصر الصفحة
    console.log('🔍 فحص عناصر صفحة تسجيل الدخول...');
    const loginFormExists = await page.$('form') !== null;
    const emailInputExists = await page.$('input[type="email"]') !== null;
    const passwordInputExists = await page.$('input[type="password"]') !== null;
    const submitButtonExists = await page.$('button[type="submit"]') !== null;
    
    console.log(`   نموذج تسجيل الدخول: ${loginFormExists ? '✅' : '❌'}`);
    console.log(`   حقل البريد الإلكتروني: ${emailInputExists ? '✅' : '❌'}`);
    console.log(`   حقل كلمة المرور: ${passwordInputExists ? '✅' : '❌'}`);
    console.log(`   زر الإرسال: ${submitButtonExists ? '✅' : '❌'}`);

    // 2. ملء نموذج تسجيل الدخول
    console.log('\n📝 ملء نموذج تسجيل الدخول...');
    await page.type('input[type="email"]', 'admin@worktrack.com');
    await page.type('input[type="password"]', 'adminadmin');
    console.log('   ✅ تم إدخال البيانات');

    // 3. إرسال النموذج
    console.log('\n🔐 إرسال نموذج تسجيل الدخول...');
    await Promise.all([
      page.click('button[type="submit"]'),
      page.waitForNavigation({ waitUntil: 'networkidle2', timeout: 10000 })
    ]);
    console.log('   ✅ تم إرسال النموذج');

    // 4. التحقق من التوجيه للوحة التحكم
    console.log('\n🎯 التحقق من التوجيه للوحة التحكم...');
    const currentUrl = page.url();
    console.log(`   الرابط الحالي: ${currentUrl}`);
    
    if (currentUrl.includes('/dashboard')) {
      console.log('   ✅ تم التوجيه للوحة التحكم بنجاح');
    } else {
      console.log('   ❌ لم يتم التوجيه للوحة التحكم');
    }

    // 5. فحص عناصر لوحة التحكم
    console.log('\n🔍 فحص عناصر لوحة التحكم...');
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    const sidebarExists = await page.$('.sidebar') !== null;
    const topbarExists = await page.$('.topbar') !== null;
    const contentExists = await page.$('.content') !== null;
    
    console.log(`   الشريط الجانبي: ${sidebarExists ? '✅' : '❌'}`);
    console.log(`   الشريط العلوي: ${topbarExists ? '✅' : '❌'}`);
    console.log(`   المحتوى الرئيسي: ${contentExists ? '✅' : '❌'}`);

    // 6. التجول في الصفحات المختلفة
    console.log('\n🧭 التجول في الصفحات المختلفة...');
    
    const pages = [
      { name: 'الموظفين', path: '/employees' },
      { name: 'نقاط العمل', path: '/worksites' },
      { name: 'التقارير', path: '/reports' },
      { name: 'إدارة الحضور', path: '/attendance-management' },
      { name: 'الملاحظات', path: '/notes' },
      { name: 'الإعدادات', path: '/settings' }
    ];

    for (const pageInfo of pages) {
      console.log(`\n   📄 الانتقال إلى ${pageInfo.name}...`);
      try {
        await page.goto(`http://localhost:3001${pageInfo.path}`, { waitUntil: 'networkidle2', timeout: 10000 });
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        const pageContent = await page.content();
        const hasContent = pageContent.length > 1000;
        
        console.log(`      ✅ تم الوصول لصفحة ${pageInfo.name}`);
        console.log(`      المحتوى: ${hasContent ? '✅ موجود' : '❌ غير موجود'}`);
      } catch (error) {
        console.log(`      ❌ خطأ في الوصول لصفحة ${pageInfo.name}: ${error.message}`);
      }
    }

    // 7. العودة للوحة التحكم
    console.log('\n🏠 العودة للوحة التحكم...');
    await page.goto('http://localhost:3001/dashboard', { waitUntil: 'networkidle2' });
    await new Promise(resolve => setTimeout(resolve, 2000));

    // 8. التحقق من وجود البيانات
    console.log('\n📊 التحقق من وجود البيانات...');
    const statsCards = await page.$$('.stat-card');
    console.log(`   بطاقات الإحصائيات: ${statsCards.length > 0 ? `✅ (${statsCards.length})` : '❌'}`);

    // 9. التحقق من الخريطة
    console.log('\n🗺️ التحقق من الخريطة...');
    const mapExists = await page.$('.dashboard__map') !== null;
    console.log(`   عنصر الخريطة: ${mapExists ? '✅' : '❌'}`);

    // 10. فحص localStorage
    console.log('\n💾 فحص localStorage...');
    const localStorageData = await page.evaluate(() => {
      return {
        user: localStorage.getItem('worktrack_admin_user'),
        theme: localStorage.getItem('worktrack_theme')
      };
    });
    console.log(`   بيانات المستخدم: ${localStorageData.user ? '✅' : '❌'}`);
    console.log(`   السمة: ${localStorageData.theme || 'غير محدد'}`);

    // 11. فحص الكوكيز
    console.log('\n🍪 فحص الكوكيز...');
    const cookies = await page.cookies();
    const hasAccessToken = cookies.some(c => c.name === 'access_token');
    const hasRefreshToken = cookies.some(c => c.name === 'refresh_token');
    console.log(`   access_token: ${hasAccessToken ? '✅' : '❌'}`);
    console.log(`   refresh_token: ${hasRefreshToken ? '✅' : '❌'}`);

    console.log('\n✅ تم الانتهاء من الاختبار بنجاح!');
    console.log('⏱️ سيتم إغلاق المتصفح خلال 5 ثواني...');
    await new Promise(resolve => setTimeout(resolve, 5000));

  } catch (error) {
    console.error('\n❌ حدث خطأ أثناء الاختبار:', error.message);
  } finally {
    await browser.close();
    console.log('🔒 تم إغلاق المتصفح');
  }
}

// تشغيل الاختبار
testApp().catch(console.error);