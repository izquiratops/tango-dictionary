class WordList {
  constructor(storageKey = 'recentWords') {
    this.storageKey = storageKey;
    this.maxItems = 8;
  }

  getWords() {
    const stored = localStorage.getItem(this.storageKey);
    return stored ? JSON.parse(stored) : [];
  }

  addWord(word) {
    let words = this.getWords();
    words = words.filter(w => w !== word);
    words.unshift(word);
    words = words.slice(0, this.maxItems);
    localStorage.setItem(this.storageKey, JSON.stringify(words));
    return words;
  }

  clearWords() {
    localStorage.removeItem(this.storageKey);
  }
}

document.addEventListener('DOMContentLoaded', () => {
  const wordList = new WordList();
  const words = wordList.getWords();

  // Render recent words list
  const recentWordsEl = document.querySelector('#recent-words')
  if (recentWordsEl) {
    recentWordsEl.innerHTML = words.map(word => `<li><a href="/search?query=${word}">${word}</a></li>`).join('');
  }

  // Add word to list on new search
  document.querySelector('form[action="/search"]').addEventListener('submit', (event) => {
    wordList.addWord(event.target.querySelector('input').value);
  });
});
