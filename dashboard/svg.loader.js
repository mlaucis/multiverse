module.exports = function (content) {
  return (
    'module.exports = require("react").createElement("div", {' +
      'className: "icon",' +
      'dangerouslySetInnerHTML: {' +
        '__html: ' + JSON.stringify(content) +
      '}' +
    '});'
  );
};
