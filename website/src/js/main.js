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

  var $demoForm = $('form#demoForm');
  var $demoSuccess = $demoForm.find('.success-feedback');

  $demoForm.on('submit', function(ev) {
    ev.preventDefault();
    ev.stopPropagation();

    var demoTrackProps = {
      firstName: $demoForm.find('#firstName').val(),
      lastName: $demoForm.find('#lastName').val(),
      email: $demoForm.find('#email').val(),
      phone: $demoForm.find('#phone').val(),
      company: $demoForm.find('#company').val(),
      comment: $demoForm.find('#comment').val(),
      companySize: $demoForm.find('#companySize').val(),
      userbase: $demoForm.find('#appSize').val()
    }

    var demoIdentifyProps = {
      firstName: $demoForm.find('#firstName').val(),
      lastName: $demoForm.find('#lastName').val(),
      email: $demoForm.find('#email').val(),
      company: $demoForm.find('#company').val()
    }

    analytics.identify(demoIdentifyProps.email, demoIdentifyProps);
    analytics.track('Demo requested', demoTrackProps);

    $demoForm.find('.uniform').slideUp(360);
    $demoSuccess.slideDown(360);
  });

  var $contentForm = $('form#contentForm');
  var $contentSuccess = $contentForm.find('.success-feedback');

  $contentForm.on('submit', function(ev) {
    ev.preventDefault();
    ev.stopPropagation();

    var contentTrackProps = {
      firstName: $contentForm.find('#firstName').val(),
      lastName: $contentForm.find('#lastName').val(),
      email: $contentForm.find('#email').val(),
      phone: $contentForm.find('#phone').val(),
      company: $contentForm.find('#company').val(),
      companySize: $contentForm.find('#companySize').val(),
      userbase: $contentForm.find('#appSize').val(),
      title: $contentForm.find('#title').val(),
      type: $contentForm.find('#type').val()
    }

    var contentIdentifyProps = {
      firstName: $contentForm.find('#firstName').val(),
      lastName: $contentForm.find('#lastName').val(),
      email: $contentForm.find('#email').val(),
      company: $contentForm.find('#company').val()
    }

    var contentTitle = $contentForm.find('#title').val();
    var requestedString = " requested";
    var eventName = contentTitle.concat(requestedString);

    analytics.identify(contentIdentifyProps.email, contentIdentifyProps);
    analytics.track('Content requested', contentTrackProps);
    analytics.track(eventName, contentTrackProps);

    $contentForm.find('.uniform').slideUp(360);
    $contentSuccess.slideDown(360);
  });

})(jQuery);
