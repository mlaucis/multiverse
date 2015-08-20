jQuery(document).ready(function($){
	/*******************
		color swatch
	********************/
	//convert rgba color to hex color
	$.cssHooks.backgroundColor = {
	    get: function(elem) {
	        if (elem.currentStyle)
	            var bg = elem.currentStyle["background-color"];
	        else if (window.getComputedStyle)
	            var bg = document.defaultView.getComputedStyle(elem,
	                null).getPropertyValue("background-color");
	        if (bg.search("rgb") == -1)
	            return bg;
	        else {
	            bg = bg.match(/^rgb\((\d+),\s*(\d+),\s*(\d+)\)$/);
	            function hex(x) {
	                return ("0" + parseInt(x).toString(16)).slice(-2);
	            }
	            return "#" + hex(bg[1]) + hex(bg[2]) + hex(bg[3]);
	        }
	    }
	}
	//set a label for each color swatch
	$('.cd-color-swatch').each(function(){
		var actual = $(this);
		$('<b>'+actual.css("background-color")+'</b>').insertAfter(actual);
	});

	/*******************
		buttons
	********************/
	$('#buttons .cd-box').each(function (idx, node) {
		console.log(arguments)
		var $el = $(node);
		var html = $el.html();
		var $container = $('<div class="cd-box"></div>').insertAfter($el);
		var text = html.split('</button>');

		$.map(text, function (value) {
			if (value.indexOf('button') >= 0) {
				var split = value.split('class="');
				var block1 = split[0] + 'clas="';
				var block2 = split[1].split('"');
				var $wrap = $('<p></p>').text(block1);
				var $span = $('<span></span>').text(block2[0]);

				$span.appendTo($wrap);
				$wrap.appendTo($container);
				$wrap.append('"' + block2[1] + '&lt;/button&gt;');
			}
		});
	});

	/*******************
		typography
	********************/
	[
		$("#typography h1"),
		$("#typography .cd-box h2"),
		$("#typography h3"),
		$("#typography p")
	].forEach(function ($el) {
		var text = $el.children('span').eq(0);

		setTypography($el, text);

		$(window).on('resize', function () {
			setTypography($el, text);
		});
	});

	function setTypography(element, textElement) {
		var fontSize = Math.round(element.css('font-size').replace('px',''))+'px',
			fontFamily = (element.css('font-family').split(','))[0].replace(/\'/g, '').replace(/\"/g, ''),
			fontWeight = element.css('font-weight');
		textElement.text(fontWeight + ' '+ fontFamily+' '+fontSize );
	}

	/*******************
		main  navigation
	********************/
	var contentSections = $('main section');
	//open navigation on mobile
	$('.cd-nav-trigger').on('click', function(){
		$('header').toggleClass('nav-is-visible');
	});
	//smooth scroll to the selected section
	$('.cd-main-nav a[href^="#"]').on('click', function(event){
      event.preventDefault();
      $('header').removeClass('nav-is-visible');
      var target= $(this.hash),
      	topMargin = target.css('marginTop').replace('px', ''),
      	hedearHeight = $('header').height();
      $('body,html').animate({'scrollTop': parseInt(target.offset().top - hedearHeight - topMargin)}, 200);
  });

  // update selected navigation element
  $(window).on('scroll', function(){
  	updateNavigation();
  });

  function updateNavigation() {
		contentSections.each(function(){
			var actual = $(this);
			var actualHeight = actual.height();
			var topMargin = actual.css('marginTop').replace('px', '');
			var actualAnchor = $('.cd-main-nav').find('a[href="#'+actual.attr('id')+'"]');

			if ( ( parseInt(actual.offset().top - $('.cd-main-nav').height() - topMargin )<= $(window).scrollTop() ) && ( parseInt(actual.offset().top +  actualHeight - (2 * topMargin) )  > $(window).scrollTop() +1 ) ) {
				actualAnchor.addClass('selected');
			} else {
				actualAnchor.removeClass('selected');
			}
		});
	}
});
