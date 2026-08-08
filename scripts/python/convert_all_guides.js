const fs = require('fs');
const path = require('path');
const { marked } = require('marked');

// Language-specific configurations
const languages = {
    arabic: {
        markdown: 'دليل_المستخدم_الشامل.md',
        html: 'دليل_المستخدم_الشامل.html',
        lang: 'ar',
        dir: 'rtl',
        title: 'دليل المستخدم الشامل لنظام WorkTrack',
        fontFamily: "'Arial', 'Tahoma', 'Times New Roman', serif"
    },
    hebrew: {
        markdown: 'מדריך_המשתמש_המלא.md',
        html: 'מדריך_המשתמש_המלא.html',
        lang: 'he',
        dir: 'rtl',
        title: 'מדריך המשתמש המלא למערכת WorkTrack',
        fontFamily: "'Arial', 'Tahoma', 'David', 'Courier New', serif"
    },
    english: {
        markdown: 'Complete_User_Guide.md',
        html: 'Complete_User_Guide.html',
        lang: 'en',
        dir: 'ltr',
        title: 'Complete User Guide for WorkTrack System',
        fontFamily: "'Arial', 'Helvetica', 'Times New Roman', serif"
    }
};

function createHTMLContent(config, htmlContent) {
    const textAlign = config.dir === 'rtl' ? 'right' : 'left';
    const paddingStart = config.dir === 'rtl' ? 'padding-right' : 'padding-left';
    
    return `<!DOCTYPE html>
<html lang="${config.lang}" dir="${config.dir}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>${config.title}</title>
    <style>
        * {
            box-sizing: border-box;
        }
        
        body {
            font-family: ${config.fontFamily};
            line-height: 1.8;
            max-width: 900px;
            margin: 2cm auto;
            padding: 2cm;
            direction: ${config.dir};
            text-align: ${textAlign};
            background-color: #ffffff;
            color: #333333;
        }
        
        h1, h2, h3, h4, h5, h6 {
            color: #2c3e50;
            margin-top: 1.5em;
            margin-bottom: 0.8em;
            font-weight: 600;
            line-height: 1.3;
        }
        
        h1 {
            font-size: 2.2em;
            border-bottom: 3px solid #3498db;
            padding-bottom: 0.4em;
            margin-top: 0;
        }
        
        h2 {
            font-size: 1.8em;
            border-bottom: 2px solid #ecf0f1;
            padding-bottom: 0.3em;
            page-break-after: avoid;
        }
        
        h3 {
            font-size: 1.4em;
            page-break-after: avoid;
        }
        
        h4 {
            font-size: 1.2em;
        }
        
        p {
            margin: 1em 0;
            text-align: justify;
        }
        
        a {
            color: #3498db;
            text-decoration: none;
            border-bottom: 1px dotted #3498db;
        }
        
        a:hover {
            color: #2980b9;
            border-bottom-style: solid;
        }
        
        code {
            background-color: #f8f9fa;
            padding: 0.2em 0.4em;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
            font-size: 0.9em;
            direction: ltr;
            display: inline-block;
            border: 1px solid #e9ecef;
        }
        
        pre {
            background-color: #f8f9fa;
            padding: 1em;
            border-radius: 5px;
            overflow-x: auto;
            direction: ltr;
            text-align: left;
            border: 1px solid #e9ecef;
            margin: 1em 0;
        }
        
        pre code {
            background-color: transparent;
            padding: 0;
            border: none;
            font-size: 0.9em;
        }
        
        blockquote {
            border-${config.dir === 'rtl' ? 'right' : 'left'}: 4px solid #3498db;
            margin: 1.5em 0;
            padding: 1em 1.5em;
            background-color: #f8f9fa;
            color: #555;
            font-style: italic;
        }
        
        table {
            border-collapse: collapse;
            width: 100%;
            margin: 1.5em 0;
            font-size: 0.95em;
        }
        
        th, td {
            border: 1px solid #ddd;
            padding: 12px;
            text-align: ${textAlign};
        }
        
        th {
            background-color: #3498db;
            color: white;
            font-weight: 600;
        }
        
        tr:nth-child(even) {
            background-color: #f8f9fa;
        }
        
        ul, ol {
            margin: 1em 0;
            ${paddingStart}: 2em;
        }
        
        li {
            margin: 0.5em 0;
            line-height: 1.6;
        }
        
        hr {
            border: none;
            border-top: 2px solid #ecf0f1;
            margin: 2em 0;
        }
        
        img {
            max-width: 100%;
            height: auto;
            margin: 1em 0;
        }
        
        /* Language-specific styling for mixed content */
        .ltr-text {
            direction: ltr;
            text-align: left;
        }
        
        .rtl-text {
            direction: rtl;
            text-align: right;
        }
        
        /* Print optimization */
        @media print {
            body {
                margin: 0;
                padding: 1cm;
            }
            
            h1 {
                page-break-before: always;
            }
            
            h2, h3 {
                page-break-after: avoid;
            }
            
            table, pre, blockquote {
                page-break-inside: avoid;
            }
            
            a {
                color: #333;
                text-decoration: none;
                border-bottom: 1px solid #333;
            }
        }
        
        /* Page numbering for PDF */
        @page {
            margin: 2cm;
        }
        
        /* Remove file paths from links in PDF */
        a[href^="file://"] {
            display: none !important;
        }
        
        /* Hide any elements that might contain file paths */
        a[href*="/home/"],
        a[href*="/Users/"],
        a[href*="C:"] {
            display: none !important;
        }
        
        /* Hide footer elements that might contain file paths */
        @page {
            @bottom-right {
                content: none;
            }
            @bottom-left {
                content: none;
            }
        }
        
        /* Remove any potential URL printing */
        @media print {
            a[href]:after {
                content: none !important;
            }
        }
    </style>
</head>
<body>
    ${htmlContent}
</body>
</html>`;
}

function removeFilePathsFromMarkdown(markdown) {
    // Remove any links containing file paths
    let cleaned = markdown.replace(/\[([^\]]+)\]\(file:\/\/[^\)]+\)/g, '$1');
    cleaned = cleaned.replace(/\[([^\]]+)\]\([^)]*\/home\/[^)]*\)/g, '$1');
    cleaned = cleaned.replace(/\[([^\]]+)\]\([^)]*\/Users\/[^)]*\)/g, '$1');
    cleaned = cleaned.replace(/\[([^\]]+)\]\([^)]*C:\\[^)]*\)/g, '$1');
    
    // Remove any inline file paths
    cleaned = cleaned.replace(/file:\/\/[^\s\)]+/g, '');
    cleaned = cleaned.replace(/\/home\/[^\s\)]+/g, '');
    cleaned = cleaned.replace(/\/Users\/[^\s\)]+/g, '');
    cleaned = cleaned.replace(/C:\\[^\s\)]+/g, '');
    
    return cleaned;
}

function convertMarkdownToHTML(config) {
    try {
        console.log(`Converting ${config.lang} guide...`);
        
        // Read markdown file
        if (!fs.existsSync(config.markdown)) {
            console.error(`❌ Markdown file not found: ${config.markdown}`);
            return false;
        }
        
        let markdownContent = fs.readFileSync(config.markdown, 'utf8');
        
        // Remove file paths from markdown content
        markdownContent = removeFilePathsFromMarkdown(markdownContent);
        
        // Configure marked options
        marked.setOptions({
            breaks: true,
            gfm: true,
            headerIds: true,
            mangle: false
        });
        
        // Convert markdown to HTML
        const htmlContent = marked(markdownContent);
        
        // Create complete HTML document
        const fullHtml = createHTMLContent(config, htmlContent);
        
        // Write HTML file
        fs.writeFileSync(config.html, fullHtml, 'utf8');
        
        console.log(`✅ Successfully created: ${config.html}`);
        return true;
        
    } catch (error) {
        console.error(`❌ Error converting ${config.lang}:`, error.message);
        return false;
    }
}

function main() {
    console.log('='*60);
    console.log('MARKDOWN TO HTML CONVERSION FOR ALL LANGUAGES');
    console.log('='*60);
    console.log();
    
    const results = {};
    
    for (const [lang, config] of Object.entries(languages)) {
        const success = convertMarkdownToHTML(config);
        results[lang] = success;
        console.log();
    }
    
    // Summary
    console.log('='*60);
    console.log('CONVERSION SUMMARY');
    console.log('='*60);
    
    const successful = Object.values(results).filter(r => r).length;
    const total = Object.keys(results).length;
    
    console.log(`Successfully converted: ${successful}/${total}`);
    console.log();
    
    for (const [lang, success] of Object.entries(results)) {
        const status = success ? '✅' : '❌';
        const config = languages[lang];
        console.log(`${status} ${lang.charAt(0).toUpperCase() + lang.slice(1)}: ${config.html}`);
    }
    
    console.log();
    console.log('Next step: Convert HTML files to PDF using:');
    console.log('  python3 convert_html_to_pdf.py');
    console.log('  or');
    console.log('  python3 simple_pdf_converter.py');
    console.log('='*60);
}

if (require.main === module) {
    main();
}

module.exports = { createHTMLContent, convertMarkdownToHTML };