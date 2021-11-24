const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    Spinner.show();
    Controller.cleanResults();
    if ((data.query || "").length === 0) {
      return;
    }
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateSearchResult(results);
        Spinner.hide();
      });
    });
  },

  getBookByTitle: (id) => {
    const form = document.getElementById(id);
    Spinner.show();
    Controller.cleanResults();
    const response = fetch(`/book?title=${form.title}`).then((response) => {
      response.json().then((results) => {
        Controller.updateBookContent(results);
        Spinner.hide();
      });
    });
  },

  updateSearchResult: (results) => {
    const booksPlace = document.getElementById("books_place");
    const quotesPlace = document.getElementById("quotes_place");
    const noResults = document.getElementById("no_results");
    const booksCards = [];
    const quotesCards = [];

    if ((results.Books || []).length > 0) {
      for (let i = 0; i < results.Books.length; i++) {
        const b = results.Books[i]
        const content = b.Chapters[0].Content
        booksCards.push(`
          <div class="card">
            <div class="card-header">
              ${b.Title}
            </div>
            <div class="card-body">
              <blockquote class="blockquote mb-0">
                <pre>${content.substr(0, 200) + "..."}</pre>
                <footer class="blockquote-footer">${b.Chapters[0].Name} <cite title="Book Name">${b.Title}</cite></footer>
              </blockquote>
              <button class="card-link btn btn-secondary get-book" id="get-book-${i}" title="${b.Title}" onclick="Controller.getBookByTitle('get-book-${i}')">Read Now</button>
            </div>
          </div>
        `);
      }
    }

    if ((results.Quotes || []).length > 0) {
      quotesCards.push(`<div class="list-group">`);
      for (let q of results.Quotes) {
        quotesCards.push(`
            <a href="#" class="list-group-item list-group-item-action flex-column align-items-start" 
                data-toggle="modal" data-target="#quote-modal" data-content="${q.Chapters[0].Content}">
              <div class="d-flex w-100 justify-content-between">
                <h5 class="mb-1">${q.Title}</h5>
              </div>
              <p class="mb-1">${q.Chapters[0].Content.substr(0, 200) + "..."}</p>
            </a>`)
        for (let it = 1; it < 5; it++) {
          if (!q.Chapters[it])
            break
          quotesCards.push(`
            <a href="#" class="list-group-item list-group-item-action flex-column align-items-start" data-toggle="modal" data-target="#quote-modal" data-content="${q.Chapters[it].Content}">
              <small class="text-muted">${q.Chapters[it].Content}</small>
            </a>
          `)
        }
        quotesCards.push(`
            <a href="#" class="list-group-item list-group-item-action flex-column align-items-start read-more" data-title="${q.Title}">
              <small class="text-muted">Read more...</small>
            </a>
          `)
      }
      quotesCards.push(`</div>`);
    }
    if (booksCards.length > 0) {
      booksPlace.innerHTML = booksCards.join('');
      booksPlace.parentElement.style.visibility = "visible";
      noResults.style.visibility = "hidden";
    }
    if (quotesCards.length > 0) {
      quotesPlace.innerHTML = quotesCards.join('');
      quotesPlace.parentElement.style.visibility = "visible";
      noResults.style.visibility = "hidden";
      bindReadMore();
    }
  },

  cleanResults: () => {
    const booksPlace = document.getElementById("books_place");
    const quotesPlace = document.getElementById("quotes_place");
    booksPlace.innerHTML = '';
    quotesPlace.innerHTML = '';
    booksPlace.parentElement.style.visibility = "hidden";
    quotesPlace.parentElement.style.visibility = "hidden";
  },

  updateBookContent: (results) => {
    const booksPlace = document.getElementById("books_place");
    const noResults = document.getElementById("no_results");
    const booksCards = [];

    if ((results || []).length > 0) {
      for (let b of results) {
        booksCards.push(`<h5 class="card-title">${b.Title}</h5>`);
        for (let c of b.Chapters) {
          booksCards.push(`
            <div class="card">
              <div class="card-body">
                <pre class="card-text">${c.Content}</pre>
              </div>
            </div>
          `)
        }
      }
    }

    booksPlace.innerHTML = booksCards.join('');
    booksPlace.parentElement.style.visibility = "visible";
    noResults.style.visibility = "hidden";
  }
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);

const buttons = document.getElementsByClassName("get-book");
for (let button of buttons) {
  button.addEventListener("click", Controller.getBookByTitle);
}

function bindReadMore() {
  $('.read-more').on('click', function (event) {
    const component = $(this)
    const title = component.data('title')
    $('#query').val(title)
    $('#search-button').trigger('click')
  });
}

$(document).ready(function(){
  Spinner();
  Spinner.hide();
  $('#quote-modal').on('show.bs.modal', function (event) {
    const component = $(event.relatedTarget)
    const quoteContent = component.data('content')
    const modal = $(this)
    modal.find('.modal-body pre').text(quoteContent)
  });
  bindReadMore();
});