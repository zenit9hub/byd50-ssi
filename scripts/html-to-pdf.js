const fs = require("fs");
const path = require("path");

const distDir = path.join(__dirname, "..", "api-docs");
const htmlPath = path.join(distDir, "api-docs.html");
const pdfPath = path.join(distDir, "api-docs.pdf");

if (!fs.existsSync(htmlPath)) {
  console.error("[html-to-pdf] missing api-docs/api-docs.html");
  process.exit(1);
}

async function run() {
  let playwright;
  try {
    playwright = require("playwright");
  } catch (err) {
    console.error("[html-to-pdf] playwright not found. Install with: npm i -D playwright");
    process.exit(1);
  }

  const browser = await playwright.chromium.launch();
  const page = await browser.newPage();
  await page.goto(`file://${htmlPath}`, { waitUntil: "networkidle" });
  await page.pdf({ path: pdfPath, format: "A4", printBackground: true });
  await browser.close();
  console.log(`[html-to-pdf] wrote ${pdfPath}`);
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
