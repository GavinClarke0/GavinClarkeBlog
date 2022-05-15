import Markdoc from '@markdoc/markdoc';
import * as fs from 'fs';

const pagesPath = './pages'
const htmlPath = './server/html'

function stripFileExtension(fileName){
    return fileName.replace(/\.[^/.]+$/, "")
}

const dir = fs.opendirSync(pagesPath)
let directFile
while ((directFile = dir.readSync()) !== null) {
    // read file content as string note file size must be at least 1/2 memory
    const content = fs.readFileSync(`${pagesPath}/${directFile.name}`).toString();

    // generate html
    const ast = Markdoc.parse(content);
    const markDownContent = Markdoc.transform(ast);
    const html = Markdoc.renderers.html(markDownContent);

    const htmlName = `${htmlPath}/${stripFileExtension(directFile.name)}.html`
    fs.writeFile(htmlName, html, function (err) {
        if (err) return console.log(err);
    });
}



