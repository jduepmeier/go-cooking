

function getCard(id) {
    let card = document.querySelector(`#recipe-${id}`);
    return card;
}

function printerhide(id) {
    let card = getCard(id);
    card.classList.add("printerhide");
    setcard(card, '.noprintericon', '.printericon')
}

function printershow(id) {
    let card = getCard(id);
    card.classList.remove("printerhide");
    setcard(card, '.printericon', '.noprintericon')
}

function setcard(card, show, hide) {
    let printericons = card.querySelector('.printericons');
    printericons.querySelector(show).classList.remove("hidden");
    printericons.querySelector(hide).classList.add("hidden");
}

function hideall() {
    let cards = document.querySelectorAll('.recipe');
    document.querySelector('.printericonsall .noprintericon').classList.remove('hidden');
    document.querySelector('.printericonsall .printericon').classList.add('hidden');
    cards.forEach((card) => {
        card.classList.add("printerhide");
        setcard(card, '.noprintericon', '.printericon');
    })
}

function showall() {
    let cards = document.querySelectorAll('.recipe');
    document.querySelector('.printericonsall .printericon').classList.remove('hidden');
    document.querySelector('.printericonsall .noprintericon').classList.add('hidden');
    cards.forEach((card) => {
        card.classList.remove("printerhide");
        setcard(card, '.printericon', '.noprintericon');
    })
}