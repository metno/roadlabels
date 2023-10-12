(function () {

    const quotesEl = document.querySelector('.quotes');
    const loaderEl = document.querySelector('.loader');

    // get the quotes from API
    const getQuotes = async (page, limit) => {
        const urlParams = new URLSearchParams(window.location.search);
        label = urlParams.get('label');
        if (label == null) {
            label = "";
        }
        const API_URL =  window.location.protocol + "//" + window.location.host + `/roadlabels/iceimagepagingapi?page=${page}&limit=${limit}&label=${label}`;
        const response = await fetch(API_URL);
        // handle 404
        if (!response.ok) {
            throw new Error(`An error occurred: ${response.status}`);
        }
        
        return await response.json();
    }

    // show the quotes
    const showQuotes = (elements) => {

        count = 0;
        tot = elements.length;

        elements.forEach(elm => {
            count++;

            const quoteEl = document.createElement('blockquote');
            quoteEl.classList.add('quote');

            if (count == 0 ) {
                quoteEl.innerHTML = '<div class="row">';
            }
            quoteEl.innerHTML += `<div class="column"><a href="/roadlabels/inputlabel?q=${elm.PathBig}"> <figure><img style="width:90%" src="/roadlabels/labeledthumb?q=${elm.PathThumb}&cc=${elm.Label}&obs2=-1"></img> <figcaption>${elm.Desc}:&nbsp;${elm.Value}</figcaption> </figure></a></div>`;

            if ( count %3 == 0) {
                quoteEl.innerHTML += '</div><div class="row">';
            }
          
            if (tot == count ) {
                console.log("count" + count + " tot: " + tot);
              
                quoteEl.innerHTML += '</div>';
            }
    
            quotesEl.appendChild(quoteEl);
        });
    };

    const hideLoader = () => {
        loaderEl.classList.remove('show');
    };

    const showLoader = () => {
        loaderEl.classList.add('show');
    };

    const hasMoreQuotes = (page, limit, total) => {
        const startIndex = (page - 1) * limit + 1;
        return total === 0 || startIndex < total;
    };

    // load quotes
    const loadQuotes = async (page, limit) => {

        // show the loader
        showLoader();

        // 0.5 second later
        setTimeout(async () => {
            try {
                // if having more quotes to fetch
                if (hasMoreQuotes(page, limit, total)) {
                    // call the API to get quotes
                    const response = await getQuotes(page, limit);
                    // show quotes
                    showQuotes(response.Data);
                    // update the total
                    total = response.Total;
                }
            } catch (error) {
                console.log(error.message);
            } finally {
                hideLoader();
            }
        }, 1000);

    };

    // control variables
    let currentPage = 1;
    const limit = 20;
    let total = 0;


    window.addEventListener('scroll', () => {
        const {
            scrollTop,
            scrollHeight,
            clientHeight
        } = document.documentElement;

        if (scrollTop + clientHeight >= scrollHeight - 5 &&
            hasMoreQuotes(currentPage, limit, total)) {
            currentPage++;
            loadQuotes(currentPage, limit);
        }
    }, {
        passive: true
    });

    // initialize
    loadQuotes(currentPage, limit);

})();