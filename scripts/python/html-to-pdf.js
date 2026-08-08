const puppeteer = require('puppeteer');
const fs = require('fs');

(async () => {
    const browser = await puppeteer.launch({
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });
    const page = await browser.newPage();
    
    // Load the HTML file
    const htmlFile = 'file://' + __dirname + '/دليل_المستخدم_الشامل.html';
    await page.goto(htmlFile, { waitUntil: 'networkidle0' });
    
    // Generate PDF
    await page.pdf({
        path: 'دليل_المستخدم_الشامل.pdf',
        format: 'A4',
        printBackground: true,
        margin: {
            top: '2cm',
            right: '2cm',
            bottom: '2cm',
            left: '2cm'
        }
    });
    
    await browser.close();
    console.log('تم إنشاء ملف PDF بنجاح: دليل_المستخدم_الشامل.pdf');
})();