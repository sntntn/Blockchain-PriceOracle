package oracle

const ContractABI = `[
	{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},
	{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"symbol","type":"string"},{"indexed":false,"internalType":"uint256","name":"oldPrice","type":"uint256"},
	{"indexed":false,"internalType":"uint256","name":"newPrice","type":"uint256"}],"name":"PriceUpdated","type":"event"},
	{"inputs":[{"internalType":"string","name":"symbol","type":"string"},{"internalType":"address","name":"tokenFeedAddress","type":"address"}],"name":"addToken","outputs":[],"stateMutability":"nonpayable","type":"function"},
	{"inputs":[{"internalType":"string","name":"","type":"string"}],"name":"chainlinkFeeds","outputs":[{"internalType":"contract AggregatorV3Interface","name":"","type":"address"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"internalType":"string","name":"symbol","type":"string"}],"name":"getPrices","outputs":[{"internalType":"int256","name":"onchainPrice","type":"int256"},{"internalType":"int256","name":"chainlinkPrice","type":"int256"}],"stateMutability":"view","type":"function"},
	{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"internalType":"string","name":"","type":"string"}],"name":"prices","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"internalType":"string","name":"symbol","type":"string"},{"internalType":"uint256","name":"newPrice","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}
]`
