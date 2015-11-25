function extractReferrer() {
  var referrer = document.referrer.split('/')[2];
  var exclude = /tapglue.com$/;

  if (!referrer || exclude.test(referrer)) {
    return
  }

  Cookies.set('originalReferrer', referrer, {
    domain: '.tapglue.com',
    expires: 7,
    path: '/'
  });
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

  startRotation();
  extractReferrer();

  var sp =new StatusPage({ pageId: '0ln51qn4551c' });

  sp.getStatus({
    success: function(data) {
      $('.status-dot').addClass(data.status.indicator);
    }
  });
})(jQuery);
