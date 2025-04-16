const path = require("path");

module.exports = {
  // 其他配置
  resolve: {
    fallback: {
      crypto: require.resolve("crypto-browserify"),
      stream: require.resolve("stream-browserify"),
    },
    external: {
      antd: "antd",
    },
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        include: [
          path.resolve(__dirname, "src"),
          path.resolve(__dirname, "node_modules/marked"),
        ],
        use: {
          loader: "babel-loader",
        },
      },
    ],
  },
};
