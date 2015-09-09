function adjustPricing($) {
  var c = Cookies.get('originalReferrer');
  var re = /producthunt.com$/;

  if (c && !re.test(c)) {
    return
  }

  $('div.tier div.price').each(function() {
    var $el = $(this);
    var price = Math.floor(parseInt($el.text()) * 0.75);

    $el.addClass('discount');
    $el.after('<div class="price">' + price + '</div>');
    $el.after('<div class="priceMonth">-25% Product Hunt Discount</div>');
  });
}

function extractReferrer() {
  var referrer = document.referrer.split('/')[2]

  if (referrer && !referrer.match('tapglue.com$')) {
    Cookies.set('originalReferrer', referrer, {
      domain: '.tapglue.com',
      expires: 7,
      path: '/'
    });
  }
}

(function($) {
  var $potraits = $('.portraits div.portrait');
  var $statements = $('.statements div.statement');
  var quoteTimer;
  var startRotation = function () {
    quoteTimer = setInterval(function () {
      $next = $potraits.filter(':not(.inactive)').next()

      if ($next.length === 0) {
        $next = $potraits.first()
      }

      $next.trigger('click');
    }, 10000);
  }

  $potraits.on('click', function (event) {
    var $el = $(this);
    var $statement = $($statements[$el.index()]);

    clearTimeout(quoteTimer);
    event.preventDefault();

    $statements.not($statement).hide();
    $statement.show();
    $potraits.not($el).addClass('inactive');
    $el.removeClass('inactive');

    startRotation();
  });

  $potraits.not(':first').addClass('inactive');
  $statements.not(':first').hide();

  adjustPricing($);
  startRotation();
  extractReferrer();
})(jQuery);
