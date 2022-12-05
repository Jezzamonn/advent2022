const fs = require('fs');
const path = require('path');

function getSumOfMispackedItemPriorities(inputText) {
    return inputText
        .split('\n')
        .map(findMispackedItem)
        .map(getItemPriority)
        .reduce((a, b) => a + b, 0);
}

function getItemPriority(item) {
    const charCode = item.charCodeAt(0);
    const aCharCode = 'a'.charCodeAt(0);
    const zCharCode = 'z'.charCodeAt(0);
    const ACharCode = 'A'.charCodeAt(0);

    if (charCode >= aCharCode && charCode <= zCharCode) {
        return charCode - aCharCode + 1;
    }
    return charCode - ACharCode + 27;
}

function findMispackedItem(rucksackStr) {
    const itemsInFirstCompartment = new Set(
        rucksackStr.slice(0, rucksackStr.length / 2));

    return Array.from(rucksackStr
        .slice(rucksackStr.length / 2))
        .find(item => itemsInFirstCompartment.has(item));
}

function findCommonItem(rucksackArr) {
    const potentialItems = new Set(rucksackArr[0]);
    for (const otherRucksack of rucksackArr.slice(1)) {
        const otherRucksackItems = new Set(otherRucksack);

        for (const item of potentialItems) {
            if (!otherRucksackItems.has(item)) {
                potentialItems.delete(item);
            }
        }
    }
    // Return the only item. Assumes input is valid :)
    return potentialItems.values().next().value;
}

function findBadges(inputText) {
    const groupSize = 3;
    const rucksacks = inputText.split('\n');

    let prioritySum = 0;

    for (let g = 0; g < rucksacks.length; g += groupSize) {
        const commonItem = findCommonItem(rucksacks.slice(g, g + groupSize));
        prioritySum += getItemPriority(commonItem);
    }
    return prioritySum;
}

function solve() {
    const filename = path.join(__dirname, 'input.txt');
    const text = fs.readFileSync(filename, 'utf-8').trim();

    const pt1 = getSumOfMispackedItemPriorities(text);
    console.log('Pt 1: ', pt1);

    const pt2 = findBadges(text);
    console.log('Pt 2: ', pt2);
}

solve();