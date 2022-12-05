const fs = require('fs');
const path = require('path');

function getMostCals(inputText) {
    let greatestCal = -1;

    const groups = inputText.trim().split('\n\n');
    for (const group of groups) {
        let calTotal = 0;
        const cals = group.split('\n');
        for (const calStr of cals) {
            const cal = parseInt(calStr);
            calTotal += cal;
        }
        if (calTotal > greatestCal) {
            greatestCal = calTotal;
        }
    }
    return greatestCal;
}

function getSumOfTopThreeCals(inputText) {
    const bestCals = [-1, -1, -1];

    const groups = inputText.trim().split('\n\n');
    for (const group of groups) {
        let calTotal = 0;
        const cals = group.split('\n');
        for (const calStr of cals) {
            const cal = parseInt(calStr);
            calTotal += cal;
        }

        if (calTotal > bestCals[2]) {
            bestCals.pop();
            bestCals.push(calTotal);
            bestCals.sort((a, b) => b - a); // Sort descending.
        }
    }

    return bestCals.reduce((a, b) => a + b, 0);
}

function solve() {
    const filename = path.join(__dirname, 'input.txt');
    const text = fs.readFileSync(filename, 'utf-8').trim();

    const result1 = getMostCals(text);
    console.log('Pt 1: ', result1);

    const result2 = getSumOfTopThreeCals(text);
    console.log('Pt 2: ', result2);
}

solve();
