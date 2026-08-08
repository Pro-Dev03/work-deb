const fs = require('fs');
const path = require('path');
const { marked } = require('marked');

// Read the markdown file
const markdownFile = 'دليل_المستخدم_الشامل.md';
const markdownContent = fs.readFileSync(markdownFile, 'utf8');

// Convert markdown to HTML
const htmlContent = marked(markdownContent);

// Create complete HTML document with proper styling for Arabic
const fullHtml = `<!DOCTYPE html>
<html lang="ar" dir="rtl">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>دليل المستخدم الشامل</title>
    <style>
        body {
            font-family: 'Arial', 'Tahoma', sans-serif;
            line-height: 1.6;
            margin: 2cm;
            direction: rtl;
            text-align: right;
        }
        h1, h2, h3, h4, h5, h6 {
            color: #333;
            margin-top: 1em;
            margin-bottom: 0.5em;
        }
        h1 {
            border-bottom: 2px solid #333;
            padding-bottom: 0.3em;
        }
        h2 {
            border-bottom: 1px solid #ccc;
            padding-bottom: 0.3em;
        }
        code {
            background-color: #f4f4f4;
            padding: 2px 4px;
            border-radius: 3px;
            font-family: monospace;
            direction: ltr;
            display: inline-block;
        }
        pre {
            background-color: #f4f4f4;
            padding: 1em;
            border-radius: 5px;
            overflow-x: auto;
            direction: ltr;
            text-align: left;
        }
        pre code {
            background-color: transparent;
            padding: 0;
        }
        blockquote {
            border-right: 4px solid #ddd;
            margin: 1em 0;
            padding: 0.5em 1em;
            background-color: #f9f9f9;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            margin: 1em 0;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: right;
        }
        th {
            background-color: #f2f2f2;
        }
        ul, ol {
            margin: 1em 0;
            padding-right: 2em;
        }
        a {
            color: #0066cc;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        @media print {
            body {
                margin: 0;
            }
            h1 {
                page-break-before: always;
            }
        }
    </style>
</head>
<body>
    ${htmlContent}
</body>
</html>`;

// Write HTML file
const htmlFile = 'دليل_المستخدم_الشامل.html';
fs.writeFileSync(htmlFile, fullHtml, 'utf8');

console.log(`تم إنشاء ملف HTML: ${htmlFile}`);
console.log('يمكنك الآن فتح الملف في المتصفح وطباعته كـ PDF باستخدام Ctrl+P');