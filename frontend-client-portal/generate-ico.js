import fs from 'fs';
import path from 'path';
import PngToIco from 'png-to-ico';

const pngPath = path.join(process.cwd(), 'public', 'icon.png');
const icoPath = path.join(process.cwd(), 'public', 'icon.ico');

PngToIco(pngPath)
  .then(buf => {
    fs.writeFileSync(icoPath, buf);
    console.log('✅ Created icon.ico successfully');
  })
  .catch(err => {
    console.error('❌ Error creating icon.ico:', err);
  });