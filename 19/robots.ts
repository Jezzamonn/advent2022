import { promises as fs } from 'fs';
import * as path from 'path';

const RESOURCE_TYPES = ['ore', 'clay', 'obsidian', 'geode'] as const;

const ZERO_RESOURCE_MAP = new Map(RESOURCE_TYPES.map((type) => [type, 0]));

class Action {
    type: string = '';
    costs: Map<string, number> = new Map();

    toString(): string {
        return `${this.type} ${this.costs}`;
    }
}

function mapToString(map: Map<string, number>): string {
    return Array.from(map.entries())
        .map(([key, value]) => `${key}=${value}`)
        .join(' ');
}

/**
 * Parses the cost of each robot type from a line of the input using regex.
 *
 * Example line:
 * `Blueprint 1: Each ore robot costs 4 ore. Each clay robot costs 4 ore. Each obsidian robot costs 4 ore and 17 clay. Each geode robot costs 4 ore and 20 obsidian.`
 */
function parseCostsToActions(line: string): Action[] {
    // First pull out the sentence that describes the costs of each robot type.
    const robotRegex = /Each (ore|clay|obsidian|geode) robot costs (.*?)\./g;
    const robotMatches = [...line.matchAll(robotRegex)];
    // Then parse the costs of each robot type.
    const actions: Action[] = [];
    for (const robotMatch of robotMatches) {
        const action = new Action();
        const robotType = robotMatch[1];
        action.type = robotType;
        const costs = robotMatch[2];
        const costRegex = /(\d+) (ore|clay|obsidian|geode)/g;
        const costMatches = [...costs.matchAll(costRegex)];
        for (const costMatch of costMatches) {
            const cost = parseInt(costMatch[1]);
            const resourceType = costMatch[2];
            action.costs.set(resourceType, cost);
        }
        actions.push(action);
    }
    return actions;
}

function maximiseGeodes(
    actions: Action[], timeLeft: number, resources: Map<string, number>,
    resourcesPerMinute: Map<string, number>, statesVisited: Map<string, number>, depth: number = 0): number {
    const depthStr = ' '.repeat(depth);
    // console.log(`${depthStr}maximiseGeodes(${timeLeft}, Resources: ${mapToString(resources)}, ResourcesPerMin: ${mapToString(resourcesPerMinute)})`);

    if (timeLeft < 0) {
        return 0;
    }

    const resourcesPerMinString = mapToString(resourcesPerMinute);

    // This might eliminate some valid solutions, but the idea is if we've had
    // the same amount of robots in a previous state, but this time we have less
    // time left, then this state must be worse than the previous state so we
    // can stop the search here.

    // Ok, is does eliminate some solutions. I guess lets just add some kludge factor here.
    const previousTimeLeft = statesVisited.get(resourcesPerMinString) ?? 0;
    if (previousTimeLeft > timeLeft + 1) {
        return 0;
    }
    else {
        statesVisited.set(resourcesPerMinString, Math.max(previousTimeLeft, timeLeft));
    }

    // How many geodes we'd have if we waited and did nothing else.
    let bestGeodes = resources.get('geode')! + timeLeft * resourcesPerMinute.get('geode')!;



    // For each type of robot, try waiting until we have enough resources and then building it.
    actionLoop: for (const action of actions) {
        // console.log(`${depthStr}Trying to build ${action.type} robot...`)
        // First wait until we have enough resources.
        let timeToWait = 0;
        const resourcesNeeded = new Map(action.costs);
        for (const [resource, amount] of resourcesNeeded) {
            const amountLeft = amount - resources.get(resource)!;
            const amountPerMinute = resourcesPerMinute.get(resource)!;
            // No amount of waiting will ever give us enough resources.
            if (amountLeft > 0 && amountPerMinute === 0) {
                // console.log(`${depthStr}Not enough ${resource} to build ${action.type} robot.`);
                continue actionLoop;
            }
            const timeToWaitForThisResource = Math.ceil(amountLeft / amountPerMinute);
            timeToWait = Math.max(timeToWait, timeToWaitForThisResource);
        }
        // If we don't have enough time left to wait, then we can't build this robot.
        if (timeToWait > timeLeft) {
            continue;
        }

        // We need to wait 1 extra minute to build the robot.
        timeToWait++;

        // Start another search after building this robot.
        const newResources = new Map();
        const newResourcesPerMinute = new Map();
        for (const resource of RESOURCE_TYPES) {
            const newAmount = resources.get(resource)! - (resourcesNeeded.get(resource) ?? 0) + timeToWait * resourcesPerMinute.get(resource)!;
            const newAmountPerMinute = resourcesPerMinute.get(resource)! + (resource === action.type ? 1 : 0);
            newResources.set(resource, newAmount);
            newResourcesPerMinute.set(resource, newAmountPerMinute);
        }

        // console.log(`${depthStr}Building ${action.type} robot after waiting ${timeToWait} minutes. ` +
            // `Resources: ${mapToString(newResources)}. Resources per minute: ${mapToString(newResourcesPerMinute)}.`);

        const geodes = maximiseGeodes(actions, timeLeft - timeToWait, newResources, newResourcesPerMinute, statesVisited, depth + 1);
        bestGeodes = Math.max(bestGeodes, geodes);
    }
    return bestGeodes;
}

async function solve(filename: string) {
    const filepath = path.join(__dirname, filename);
    const inputText = (await fs.readFile(filepath, 'utf-8')).trimEnd();

    console.log(`Solving for ${filename}...`);

    const lines = inputText.split('\n');
    let qualityNumberSum = 0;
    for (const [i, line] of lines.entries()) {
        const actions = parseCostsToActions(line);

        process.stdout.write(`Blueprint ${i + 1}: `);

        // console.log(actions);

        // Start with no resources.
        const resources = new Map(ZERO_RESOURCE_MAP);
        // But we do start with one ore-producing robot.
        const resourcesPerMinute = new Map(ZERO_RESOURCE_MAP);
        resourcesPerMinute.set('ore', 1);

        const geodes = maximiseGeodes(actions, 24, resources, resourcesPerMinute, new Map());
        const qualityNumber = geodes * (i + 1);

        qualityNumberSum += qualityNumber;

        console.log(`${geodes} geodes. Quality number: ${qualityNumber} (sum = ${qualityNumberSum}))`);
    }

    console.log(`Part 1: ${qualityNumberSum}`);
}

Promise.resolve()
.then(() => solve('demo.txt'))
// .then(() => solve('input.txt'))