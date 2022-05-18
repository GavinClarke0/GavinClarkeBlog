import Markdoc from '@markdoc/markdoc';
import * as fs from 'fs';

const pagesPath = './pages'
const htmlPath = './server/html'
const componentsPath = './components'

function stripFileExtension(fileName){
    return fileName.replace(/\.[^/.]+$/, "")
}

function formatHtml( head, header, body, footer) {
    return `<!DOCTYPE html>\n<html lang="en">\n${head}\n<body>${header}\n${body}\n${footer}\n</body></html>`
}

const header = fs.readFileSync(`${componentsPath}/header.html`).toString();
const head = fs.readFileSync(`${componentsPath}/head.html`).toString();
const footer = fs.readFileSync(`${componentsPath}/footer.html`).toString();


const dir = fs.opendirSync(pagesPath)
let directFile
while ((directFile = dir.readSync()) !== null) {
    // read file content as string note file size must be at least 1/2 memory
    const content = fs.readFileSync(`${pagesPath}/${directFile.name}`).toString();

    // generate html
    const ast = Markdoc.parse(content);
    const markDownContent = Markdoc.transform(ast);
    const bodyHtml = Markdoc.renderers.html(markDownContent);

    const htmlName = `${htmlPath}/${stripFileExtension(directFile.name)}.html`

    const html = formatHtml(head, header, bodyHtml, footer)

    fs.writeFile(htmlName, html, function (err) {
        if (err) return console.log(err);
    });
}



