// Avoid `console` errors in browsers that lack a console.
(function() {
    var method;
    var noop = function () {};
    var methods = [
        'assert', 'clear', 'count', 'debug', 'dir', 'dirxml', 'error',
        'exception', 'group', 'groupCollapsed', 'groupEnd', 'info', 'log',
        'markTimeline', 'profile', 'profileEnd', 'table', 'time', 'timeEnd',
        'timeline', 'timelineEnd', 'timeStamp', 'trace', 'warn'
    ];
    var length = methods.length;
    var console = (window.console = window.console || {});

    while (length--) {
        method = methods[length];

        // Only stub undefined methods.
        if (!console[method]) {
            console[method] = noop;
        }
    }
}());

// Place any jQuery/helper plugins in here.
        // Request Access
        function requestAccess(el) {
            $('.request-access').show();
        }

        function hideRequestAccess() {
            $('.request-access').hide();
        }

        // Send Request Event
        function requestSubmit() {
            var email = $('.request-access input[name="email"]'),
            name = $('.request-access input[name="name"]'),
            companyname = $('.request-access input[name="companyname"]');

            if (typeof analytics !== 'undefined') {
                var data = {
                    email: email.val(),
                    firstName: name.val().split(' ')[0],
                    lastName: name.val().split(' ')[1] || '',
                    referrer: document.referrer.split('/')[2],
                    companyName: companyname.val()
                };
                analytics.identify(email.val(), data);
                analytics.track('Request access submitted', data);
            }

            $('.request-access h2').html('Thanks! We received your message and will come back to you shortly.');
            $('.request-access h2').addClass("success-request");

            email.val('');
            name.val('');
            companyname.val('');
            $('.request-access input[type="submit"]').attr('disabled', true);
            setTimeout(function () {
                $('.request-access').hide();
            }, 3500);
        }

        // Submit Form
        function feedbackSubmit() {
                var name = $('#feedbackForm input[name="name"]').val(),
                email = $('#feedbackForm input[name="email"]').val(),
                message = $('#feedbackForm textarea[name="message"]').val();

                // if (typeof analytics !== 'undefined') {
                //     analytics.identify(email, {name: name, email: email, message: message});
                //     analytics.track('Form submitted', {name: name, email: email, message: message});
                // }

                $.ajax({
                  url: 'https://www.tapglue.com/submitMailer.php',
                  type: 'post',
                  dataType: 'json',
                  data: $('form#feedbackForm').serialize(),
                  success: function(data) {
                    $('#feedbackForm input[name="name"]').val('');
                    $('#feedbackForm input[name="email"]').val('');
                    $('#feedbackForm textarea[name="message"]').val('');
                    $('#feedbackForm input[type="submit"]').attr('disabled', true);
                    $('.success-feedback').html('Thanks! We received your message and will come back to you shortly.');
                  },
                error: function(xhr, ajaxOptions, thrownError) {
                  $('.success-feedback').html('Oh snap, something went wrong. Please try again.');
                }
               });
        }

         // Submit Newsletter
        function newsletterSubmit() {
                var form = $('#newsform').val();

                $('#newsform input[name="email"]').val('');
                $('#newsform input[type="submit"]').val('Sending...');
                $('#newsform input[type="submit"]').attr('disabled', true);

                form.submit();
        }

        // Tabs and Phone Screens
        (function($){
            /* trigger when page is ready */
            $(document).ready(function (){

                //Tabs functionality
                //Firstly hide all content divs
                $('#generic-tabs div').hide();
                //Then show the first content div
                $('#generic-tabs div:first').show();
                //Add active class to the first tab link
                $('#generic-tabs ul#tabs li:first').addClass('active');
                //Functionality when a tab is clicked
                $('#generic-tabs ul#tabs li a').click(function(){
                    //Firstly remove the current active class
                    $('#generic-tabs ul#tabs li').removeClass('active');
                    //Apply active class to the parent(li) of the link tag
                    $(this).parent().addClass('active');
                    //Set currentTab to this link
                    var currentTab = $(this).attr('href');
                    //Hide away all tabs
                    $('#generic-tabs div').hide();
                    //show the current tab
                    $(currentTab).show();
                    //Stop default link action from happening
                    return false;
                });

                //Phone Screen Changer
                var iphoneSelector = '.iphone-selector';
                    $(iphoneSelector).on('click', function(){
                        $(iphoneSelector).removeClass('active');
                        $(this).addClass('active');
                    });

                //Change display of our iPhone
                $( ".iphone-selector" ).click(function() {
                    $(".iphone").attr("id",this.id + "-display");
                });

                // Show navigation on scroll
                // $(window).scroll(function(){
                //     if ($(this).scrollTop() > 600) {
                //         $('.menu').slideDown(300);
                //     } else {
                //         $('.menu').slideUp(200);
                //     }
                // });

            });
        })(window.jQuery);

        $( document ).ready(function() {

    //keep our elements active
    var iphoneSelector = '.iphone-selector';
    $(iphoneSelector).on('click', function(){
        $(iphoneSelector).removeClass('active');
        $(this).addClass('active');
    });

    //change display of our iPhone
    $( ".iphone-selector" ).click(function() {
        $(".iphone").attr("id",this.id + "-display");
    });

});
