module.exports = {
  install: function (less, pluginManager, functions) {
    functions.add('px2rem', function ({ value }) {
      return `${value / 100}rem`;
    });
  }
};