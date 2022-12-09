const fs = require('fs');
const path = require('path');

interface FileOrDirectory {
    name: string;
    size: number;

    parent?: Directory;

    // Iterate over the contents of this directory or file
    iterContents(): Iterable<FileOrDirectory>;

    print(depth: number): string;
}

class File_ implements FileOrDirectory {
    name: string;
    size: number;
    parent?: Directory;

    constructor(name: string, size: number) {
        this.name = name;
        this.size = size;
    }

    *iterContents(): Iterable<FileOrDirectory> {
        yield this;
    }

    print(depth: number) {
        const depthStr = '  '.repeat(depth);
        return `${depthStr}- ${this.name} ${this.size.toLocaleString()}\n`;
    }
}

class Directory implements FileOrDirectory {
    name: string;
    contents: Map<string, FileOrDirectory> = new Map();
    parent?: Directory;

    // Cache size to make iterating through everything not as insane.
    #cachedSize: number | undefined;

    constructor(name: string) {
        this.name = name;
    }

    get size(): number {
        if (this.#cachedSize === undefined) {
            let sum = 0;
            for (const fileOrDirectory of this.contents.values()) {
                sum += fileOrDirectory.size;
            }
            this.#cachedSize = sum;
        }
        return this.#cachedSize;
    }

    *iterContents(): Iterable<FileOrDirectory> {
        yield this;
        for (const fileOrDirectory of this.contents.values()) {
            yield* fileOrDirectory.iterContents();
        }
    }

    print(depth: number) {
        const depthStr = '  '.repeat(depth);

        let str = `${depthStr}v ${this.name}: (${this.size.toLocaleString()})\n`;

        for (const fileOrDirectory of this.contents.values()) {
            str += fileOrDirectory.print(depth + 1);
        }
        return str;
    }
}

function parseLogs(inputText: string) {
    const rootDirectory = new Directory('/');
    rootDirectory.name = '/';

    let currentDirectory = rootDirectory;

    function handleCommand(command: string[]) {
        switch (command[0]) {
            case 'cd':
                const directoryName = command[1];
                if (directoryName == '/') {
                    currentDirectory = rootDirectory;
                }
                else if (directoryName == '..') {
                    currentDirectory = currentDirectory.parent!;
                }
                else {
                    const directory = currentDirectory.contents.get(directoryName) as Directory;
                    currentDirectory = directory;
                }
                break;
            // Ignore ls commands
        }
    }

    function handleLsOutput(line: string) {
        const parts = line.split(' ');
        const name = parts[1];

        // Don't add the same file / directory twice
        if (currentDirectory.contents.has(name)) return;

        if (parts[0] == 'dir') {
            const directory = new Directory(name);
            directory.parent = currentDirectory;
            currentDirectory.contents.set(name, directory);
        }
        else {
            const size = parseInt(parts[0]);
            const file = new File_(name, size);
            file.parent = currentDirectory;
            currentDirectory.contents.set(name, file);
        }
    }

    for (const line of inputText.split('\n')) {
        if (line == '') continue;

        if (line.startsWith('$')) {
            handleCommand(line.split(' ').slice(1));
        }
        else {
            // I think the only other thing that can generate stuff is the ls
            // command. We don't really need to know that ls was even typed,
            // we can just add the files to the current directory.
            handleLsOutput(line);
        }
    }

    return rootDirectory;
}

function sumDirectorySizesLessThan(rootDirectory: Directory, size: number) {
    let sum = 0;
    for (const fileOrDirectory of rootDirectory.iterContents()) {
        // Ignore files
        if (!(fileOrDirectory instanceof Directory)) continue;

        if (fileOrDirectory.size <= size) {
            sum += fileOrDirectory.size;
        }
    }
    return sum;
}

function findSmallestDirectoryThatCanBeDeleted(rootDirectory: Directory) {
    const requiredSize = 70_000_000 - 30_000_000;

    const totalSize = rootDirectory.size;

    const sizeDifference = totalSize - requiredSize;
    let smallestDirectoryThatMakesEnoughSpace: Directory | undefined = undefined;

    for (const fileOrDirectory of rootDirectory.iterContents()) {
        // Ignore files
        if (!(fileOrDirectory instanceof Directory)) continue;

        // Not big enough to make enough space.
        if (fileOrDirectory.size < sizeDifference) continue;

        if (smallestDirectoryThatMakesEnoughSpace === undefined || fileOrDirectory.size < smallestDirectoryThatMakesEnoughSpace.size) {
            smallestDirectoryThatMakesEnoughSpace = fileOrDirectory;
        }
    }
    return smallestDirectoryThatMakesEnoughSpace?.size;
}

function solve() {
    const filename = path.join(__dirname, 'input.txt');
    const text = fs.readFileSync(filename, 'utf-8');

    const rootDirectory = parseLogs(text);
    console.log(rootDirectory.print(0));

    const pt1 = sumDirectorySizesLessThan(rootDirectory, 100_000);
    console.log('Part 1:', pt1);

    const pt2 = findSmallestDirectoryThatCanBeDeleted(rootDirectory);
    console.log('Part 2:', pt2);
}

solve();