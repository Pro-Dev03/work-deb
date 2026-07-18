import sharp from 'sharp';

async function generateIcons() {
  const inputImage = './public/company-logo.jpg';
  
  // Generate different sizes for icons
  const sizes = [16, 32, 48, 64, 128, 256, 512];
  
  for (const size of sizes) {
    await sharp(inputImage)
      .resize(size, size, { fit: 'cover', position: 'center' })
      .png()
      .toFile(`./public/icon-${size}x${size}.png`);
    
    console.log(`Generated icon-${size}x${size}.png`);
  }
  
  // Generate main icon.png (512x512)
  await sharp(inputImage)
    .resize(512, 512, { fit: 'cover', position: 'center' })
    .png()
    .toFile('./public/icon.png');
  
  console.log('Generated icon.png (512x512)');
  
  console.log('All icons generated successfully!');
}

generateIcons().catch(console.error);