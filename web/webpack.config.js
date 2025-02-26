module.exports = {
  // 其他配置
  resolve: {
    fallback: {
      crypto: require.resolve('crypto-browserify'),
      stream: require.resolve('stream-browserify'),
    },
    external: {
      'antd': 'antd',
    }
  },
};
