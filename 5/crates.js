const fs = require('fs');
const path = require('path');

const minHeight = 30;

async function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function parseAndSolveCrates(inputText, applyInstructionFn) {
    const [crates, instructions] = inputText.split('\n\n');
    const stacks = parseInitialCrates(crates);

    for (const instructionStr of instructions.split('\n')) {
        if (instructionStr === '') {
            break;
        }
        const instruction = parseInstruction(instructionStr);
        await applyInstructionFn(stacks, instruction);
    }
    return stacks.map(stack => stack[stack.length - 1]).join('');
}

async function applySingleInstructionOneBoxAtATime(stacks, instruction) {
    for (let i = 0; i < instruction.amount; i++) {
        stacks[instruction.to].push(stacks[instruction.from].pop());

        await animateCrates(stacks);
    }
}

async function applySingleInstructionGroupedBoxes(stacks, instruction) {
    const crates = stacks[instruction.from].splice(-instruction.amount);
    stacks[instruction.to].push(...crates);

    await animateCrates(stacks);
}

async function animateCrates(stacks) {
    console.clear();
    console.log(drawCrates(stacks, minHeight));
    await sleep(30);
}

function parseInitialCrates(inputText) {
    const lines = inputText.split('\n');

    const numStacks = Math.ceil(lines[0].length / 4);
    const stacks = Array(numStacks).fill(null).map(() => []);
    for (const line of lines) {
        if (line.indexOf('[') === -1) {
            break;
        }

        for (let i = 0; i < numStacks; i++) {
            const crate = line.charAt(i * 4 + 1);
            if (crate === ' ') {
                continue;
            }
            stacks[i].unshift(crate);
        }
    }
    return stacks;
}

function drawCrates(stacks, minHeight = 0) {
    const highestStack = Math.max(...stacks.map(s => s.length), minHeight);
    let output = '';
    for (let i = highestStack - 1; i >= 0; i--) {
        for (let s = 0; s < stacks.length; s++) {
            const stack = stacks[s];
            if (i < stack.length) {
                output += `[${stack[i]}]`
            }
            else {
                output += '   ';
            }

            if (s < stacks.length - 1) {
                output += ' ';
            }
        }
        output += '\n';
    }
    return output;
}

function parseInstruction(instruction) {
    // Example instruction: "move 1 from 2 to 1"
    const parts = instruction.split(' ');
    return {
        amount: parseInt(parts[1]),
        // Deal with 1-based indexing
        from: parseInt(parts[3]) - 1,
        to: parseInt(parts[5]) - 1,
    };
}

async function solve() {
    const filename = path.join(__dirname, 'input.txt');
    const text = fs.readFileSync(filename, 'utf-8');

    const pt1 = await parseAndSolveCrates(text, applySingleInstructionOneBoxAtATime);
    console.log('Part 1:', pt1);

    const pt2 = await parseAndSolveCrates(text, applySingleInstructionGroupedBoxes);
    console.log('Part 2:', pt2);
}

solve();