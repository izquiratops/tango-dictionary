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
    const words = this.getWords();
    words = words.filter(w => w !== word); // Remove existing word
    words.unshift(word); // Add new word
    words = words.slice(0, this.maxItems); // Limit to max item count
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
  document.querySelector('#recent-words').innerHTML = words.map(word => `<li>${word}</li>`).join('');

  // Add word to list on new search
  document.querySelector('form[action="/search"]').addEventListener('submit', (event) => {
    const word = event.target.querySelector('input').value;
    wordList.addWord(word);
  });
});
