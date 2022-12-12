const fs = require('fs');
const path = require('path');

// Example diagram so I understand timing
// Cycles:      |--1--|--2--|--3--|--4--|--5--|--6--|--7--|--8--|--9--|
// Instruction:  noop  addx  ----  noop  addx  ----  addx  ----  noop
// Value of x:  0     0     0     1     1     1     2     2     3     3
// cyclenum:    0     1     2     3     4     5     6     7     8     9
// Signal str:  -  0  -  0  -  0  - 4x1 - 5x1 - 6x1 - 7x2 - 8x2 - 9x3 -

// ... eh. I'll just do it the easier way.

function isImportantSignal(cycleNumber) {
    return (cycleNumber + 20) % 40 === 0;
}

function simulateInstructions(inputText) {
    let xRegister = 1;
    let cycleNumber = 0;
    let signalStrSum = 0;

    let image = '';

    function printHeader() {
        console.log('Instruction Cycle xRegister SignalStr');
    }

    function printCycle(line, cycleNumber, xRegister) {
        const signalStregth = isImportantSignal(cycleNumber) ? xRegister * cycleNumber : '-';
        // Print, while padding to make sure that the numbers are aligned.
        const linePadded = line.padEnd('Instruction'.length);
        const cycleNumberPadded = cycleNumber.toString().padStart('Cycle'.length);
        const xRegisterPadded = xRegister.toString().padStart('xRegister'.length);
        const signalStregthPadded = signalStregth.toString().padStart('SignalStr'.length);
        console.log(`${linePadded} ${cycleNumberPadded} ${xRegisterPadded} ${signalStregthPadded}`);
    }

    function addIfImportantSignal() {
        if (isImportantSignal(cycleNumber)) {
            signalStrSum += xRegister * cycleNumber;
        }
    }

    function handleCycle(line) {
        printCycle(line, cycleNumber, xRegister);
        addToImage();
        addIfImportantSignal();
    }

    function addToImage() {
        const imageWidth = 40;
        let pixelCol = (cycleNumber - 1 + imageWidth) % imageWidth;

        if (xRegister === pixelCol || xRegister === pixelCol + 1 || xRegister === pixelCol - 1) {
            image += '#';
        }
        else {
            image += ' ';
        }

        if (pixelCol === imageWidth - 1) {
            image += '\n';
        }
    }

    printHeader();
    for (const line of inputText.split('\n')) {
        instructions = line.split(' ');
        if (instructions[0] === 'noop') {
            cycleNumber++;
            handleCycle(line);
        }
        else if (instructions[0] === 'addx') {
            cycleNumber++;
            handleCycle(line);
            cycleNumber++;
            handleCycle(line);
            xRegister += parseInt(instructions[1]);
        }
    }

    console.log();
    console.log(image);

    return signalStrSum;
}

function solve(filename) {
    const filepath = path.join(__dirname, filename);
    const text = fs.readFileSync(filepath, 'utf-8').trimEnd();

    console.log(`${filename}:`);

    const pt1 = simulateInstructions(text);
    console.log('Part 1:', pt1);
}

solve('demo.txt');
solve('demo2.txt');
solve('input.txt');