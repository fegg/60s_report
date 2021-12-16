module.exports = {
  install: function (less, pluginManager, functions) {
    functions.add('px2rem', function (val) {
      return `${val / 100}rem`;
    });
  }
};