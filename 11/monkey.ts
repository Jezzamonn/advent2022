const fs = require('fs');
const path = require('path');

class Monkey {
    index: number = -1;
    items: bigint[] = [];
    operation: (item: bigint) => bigint = (item: bigint) => 0n;
    testDivisor: bigint = 0n;
    destIfTrue: number = 0;
    destIfFalse: number = 0;

    inspectCount: bigint = 0n;

    test(item: bigint): boolean {
        return item % this.testDivisor === 0n;
    }

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
        monkey.testDivisor = BigInt(linesSplit[3][3]);
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

/**
 * Print a log message like the following:
 * ```
 * == After round 20 ==
 * Monkey 0 inspected items 99 times.
 * Monkey 1 inspected items 97 times.
 * Monkey 2 inspected items 8 times.
 * Monkey 3 inspected items 103 times.
 * ```
 */
function printRoundSummary(monkeys: Monkey[], round: number) {
    console.log(`== After round ${round} ==`);
    for (const monkey of monkeys) {
        console.log(`Monkey ${monkey.index} inspected items ${monkey.inspectCount} times.`);
    }
}

function simulateMonkeys(inputText: string, rounds: number, relief = true): bigint {
    const monkeys = parseMonkeys(inputText);

    const multipleOfDivisors = monkeys.map(m => m.testDivisor).reduce((a, b) => a * b);

    for (let r = 0; r < rounds; r++) {
        for (const monkey of monkeys) {
            while (monkey.items.length > 0) {
                const item = monkey.items.shift()!;

                let newItem = monkey.operation(item);
                // I don't know if this works but why not try moduloing the thing?
                monkey.inspectCount++;

                // Experience 'relief', dividing the item worry level by 3
                if (relief) {
                    newItem = newItem / 3n;
                }
                newItem = newItem % multipleOfDivisors;

                const dest = monkey.test(newItem) ? monkey.destIfTrue : monkey.destIfFalse;
                monkeys[dest].items.push(newItem);

                // printSummaries(monkeys);
            }
        }
        const roundName = r + 1;
        if (roundName == 1 || roundName == 20 ||
            (roundName > 0 && roundName % 1000 === 0)) {
            printRoundSummary(monkeys, roundName);
        }
    }

    const inspectCounts = monkeys.map(m => m.inspectCount);
    inspectCounts.sort((a, b) => {
        if (a > b) {
            return -1;
        } else if (a < b) {
            return 1;
        } else {
            return 0;
        }
    });
    return inspectCounts[0] * inspectCounts[1];
}

function solve(filename: string) {
    const filepath = path.join(__dirname, filename);
    const inputText = fs.readFileSync(filepath, 'utf-8').trimEnd();

    console.log(`${filename}:`);
    const pt1 = simulateMonkeys(inputText, 20, true);
    console.log(`Part 1: ${pt1}`);

    const pt2 = simulateMonkeys(inputText, 10000, false);
    console.log(`Part 2: ${pt2}`);
    console.log();
}

solve('demo.txt');
solve('input.txt');
