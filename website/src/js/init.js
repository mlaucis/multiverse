/*
	Transit by TEMPLATED
	templated.co @templatedco
	Released for free under the Creative Commons Attribution 3.0 license (templated.co/license)
*/

(function($) {

	skel.init({
		reset: 'full',
		breakpoints: {
			global: {
				grid: { gutters: [ '4rem', '0' ] },
				containers: 1200
			},
			xlarge: {
				media: '(max-width: 1680px)',
				href: 'css/style-xlarge.css'
			},
			large: {
				media: '(max-width: 1280px)',
				href: 'css/style-large.css',
				viewport: { scalable: false }
			},
			medium: {
				media: '(max-width: 1024px)',
				href: 'css/style-medium.css',
				containers: '90%!'
			},
			small: {
				media: '(max-width: 768px)',
				href: 'css/style-small.css',
				containers: '90%!',
			},
			xsmall: {
				media: '(max-width: 480px)',
				href: 'css/style-xsmall.css'
			}
		},
		plugins: {
			layers: {
				config: {
					mode: 'transform'
				},
				navButton: {
					breakpoints: 'medium',
					height: '4em',
					html: '<span class="toggle" data-action="toggleLayer" data-args="navPanel"></span>',
					position: 'top-right',
					side: 'top',
					width: '6em'
				},
				navPanel: {
	        animation: 'pushX',
					breakpoints: 'medium',
					clickToHide: true,
					height: '100%',
					hidden: true,
					html: '<div data-action="moveElement"></div>',
					orientation: 'vertical',
					position: 'top-right',
					side: 'right',
					width: 200,
					html: 	'<ul>' +
							'<li><a href="/">Home</a></li>' +
							'<li><a href="/features/">Features</a></li>' +
  						'<li><a href="/news-feed/">News Feed</a></li>' +
							'<li><a href="/elements/">Elements</a></li>' +
							'<li><a href="/pricing">Pricing</a></li>' +
  						'<li><a href="//developers.tapglue.com">Docs</a></li>' +
  						'<li><a href="/blog/">Blog</a></li>' +
  						'<li><a href="/about-us/">About us</a></li>' +
  					'</ul>' +
						'<a href="/demo/" class="btn-default block navbtn">Demo</a>' +
						'<a href="//dashboard.tapglue.com/signup" class="btn-default block navbtn">Signup</a>'
				}
			}
		}
	});

	$(function() {

		var	$window = $(window),
			$body = $('body');

		// Disable animations/transitions until the page has loaded.
			$body.addClass('is-loading');

			$window.on('load', function() {
				$body.removeClass('is-loading');
			});

	});

})(jQuery);
