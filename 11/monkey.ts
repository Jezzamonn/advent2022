const fs = require('fs');
const path = require('path');

class Monkey {
    index: number = -1;
    items: number[] = [];
    operation: (item: number) => number = (item: number) => 0;
    test: (item: number) => boolean = (item: number) => false;
    destIfTrue: number = 0;
    destIfFalse: number = 0;

    inspectCount: number = 0;

    /**
     * Parses information from the string representation.
     *
     * Example string:
     * ```
     * Monkey 0:
     *  Starting items: 89, 74
     *  Operation: new = old * 5
     *  Test: divisible by 17
     *    If true: throw to monkey 4
     *    If false: throw to monkey 7
     * ```
     *
     * Thank you Copilot for this spaghetti code :)
     */
    static parse(line: string): Monkey {
        const monkey = new Monkey();

        const lines = line.split('\n');

        monkey.index = parseInt(lines[0].trim().split(' ')[1]);
        monkey.items = lines[1].trim().split(' ').slice(2).map((x: string) => parseInt(x));
        const operation = lines[2].trim().split(' ')[4];
        const operationNumber = parseInt(lines[2].trim().split(' ')[5]);

        monkey.operation = (item: number) => {
            if (isNaN(operationNumber)) {
                return item * item;
            }
            switch (operation) {
                case '+':
                    return item + operationNumber;
                case '*':
                    return item * operationNumber;
                default:
                    throw new Error(`Unknown operation ${operation}`);
            }
        }
        const testDivisor = parseInt(lines[3].trim().split(' ')[3]);
        monkey.test = (item: number) => item % testDivisor === 0;
        monkey.destIfTrue = parseInt(lines[4].trim().split(' ')[5]);
        monkey.destIfFalse = parseInt(lines[5].trim().split(' ')[5]);

        return monkey;
    }

    summary(): string {
        return `${this.index}: ${this.items.join(', ')}`;
    }
}

function printSummaries(monkeys: Monkey[]) {
    console.log(monkeys.map(m => m.summary()).join('\n') + '\n');
}

function parseMonkeys(inputText: string): Monkey[] {
    return inputText
        .split('\n\n')
        .map(Monkey.parse)
}

function simulateMonkeys(inputText: string) {
    const monkeys = parseMonkeys(inputText);

    const rounds = 20;
    for (let r = 0; r < rounds; r++) {
        for (const monkey of monkeys) {
            // Copy the items array so we can modify it while iterating
            const items = monkey.items;
            monkey.items = [];
            for (const item of items) {
                let newItem = monkey.operation(item);
                monkey.inspectCount++;
                // Experience 'relief', dividing the item worry level by 3
                newItem = Math.floor(newItem / 3);

                const dest = monkey.test(newItem) ? monkey.destIfTrue : monkey.destIfFalse;
                monkeys[dest].items.push(newItem);

                printSummaries(monkeys);
            }
        }
    }

    const inspectCounts = monkeys.map(m => m.inspectCount);
    inspectCounts.sort((a, b) => b - a);
    return inspectCounts[0] * inspectCounts[1];
}

function solve(filename: string) {
    const filepath = path.join(__dirname, filename);
    const inputText = fs.readFileSync(filepath, 'utf-8').trimEnd();

    const pt1 = simulateMonkeys(inputText);
    console.log(`Part 1: ${pt1}`);
}

solve('input.txt');
