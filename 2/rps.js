const fs = require('fs');
const path = require('path');

const opponentShapes = {
    'A': 'R',
    'B': 'P',
    'C': 'S',
}

const playerShapes = {
    'X': 'R',
    'Y': 'P',
    'Z': 'S',
}

const letterToGameResult = {
    'X': 'loss',
    'Y': 'draw',
    'Z': 'win',
}

const shapeScores = {
    'R': 1,
    'P': 2,
    'S': 3,
}

const resultScores = {
    'win': 6,
    'draw': 3,
    'loss': 0,
}

function getScoreOfRockPaperScissorsGames(inputText) {
    let totalScore = 0;

    const games = inputText.split('\n');
    for (const game of games) {
        const shapes = game.split(' ');
        const opponentShape = opponentShapes[shapes[0]];
        const playerShape = playerShapes[shapes[1]];

        const result = getResult(opponentShape, playerShape);

        const score = shapeScores[playerShape] + resultScores[result];
        totalScore += score;
    }
    return totalScore;
}

function getTotalScoreFromResults(inputText) {
    let totalScore = 0;

    const games = inputText.split('\n');
    for (const game of games) {
        const gameInfo = game.split(' ');
        const opponentShape = opponentShapes[gameInfo[0]];
        const result = letterToGameResult[gameInfo[1]];
        const playerShape = inferShapeFromResult(opponentShape, result);

        const score = shapeScores[playerShape] + resultScores[result];
        totalScore += score;
    }
    return totalScore;
}

function inferShapeFromResult(opponentShape, result) {
    // Very lazy... but simple :)
    for (const shape of ['R', 'P', 'S']) {
        if (getResult(opponentShape, shape) === result) {
            return shape;
        }
    }
}


function getResult(opponentShape, playerShape) {
    if (opponentShape === playerShape) {
        return 'draw';
    }
    if (opponentShape === 'R' && playerShape === 'P' ||
        opponentShape === 'P' && playerShape === 'S' ||
        opponentShape === 'S' && playerShape === 'R') {
        return 'win';
    }
    return 'loss';
}

function solve() {
    const filename = path.join(__dirname, 'input.txt');
    const text = fs.readFileSync(filename, 'utf-8').trim();

    const result1 = getScoreOfRockPaperScissorsGames(text);
    console.log('Pt 1: ', result1);

    const result2 = getTotalScoreFromResults(text);
    console.log('Pt 2: ', result2);
}

solve();