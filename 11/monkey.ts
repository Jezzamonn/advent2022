const fs = require('fs');
const path = require('path');

// BigInt solution is too slow. Will return and try a different approach.

class Monkey {
    index: number = -1;
    items: bigint[] = [];
    operation: (item: bigint) => bigint = (item: bigint) => 0n;
    test: (item: bigint) => boolean = (item: bigint) => false;
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

        const linesSplit = line.split('\n').map((line) => line.trim().split(' '));

        monkey.index = parseInt(linesSplit[0][1]);
        monkey.items = linesSplit[1].slice(2).map((item) => BigInt(parseInt(item)));
        const operation = linesSplit[2][4];

        const operationNumberStr = linesSplit[2][5];
        let operationNumber = 0n;
        if (operationNumberStr !== 'old') {
            operationNumber = BigInt(linesSplit[2][5]);
        }

        monkey.operation = (item) => {
            if (operationNumberStr === 'old') {
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
        const testDivisor = BigInt(linesSplit[3][3]);
        monkey.test = (item) => item % testDivisor === 0n;
        monkey.destIfTrue = parseInt(linesSplit[4][5]);
        monkey.destIfFalse = parseInt(linesSplit[5][5]);

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

    const rounds = 10000;
    for (let r = 0; r < rounds; r++) {
        for (const monkey of monkeys) {
            while (monkey.items.length > 0) {
                const item = monkey.items.shift()!;

                let newItem = monkey.operation(item);
                monkey.inspectCount++;

                // Part 2: No relief.
                // // Experience 'relief', dividing the item worry level by 3
                // newItem = newItem / 3n;

                const dest = monkey.test(newItem) ? monkey.destIfTrue : monkey.destIfFalse;
                monkeys[dest].items.push(newItem);

                // printSummaries(monkeys);
            }
        }
        if (r % 100 === 0) {
            console.log(`Round ${r} complete`);
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
