package oracle

const ContractABI = `[
	{"anonymous":false,"inputs":[{"indexed":true,"internalType":"string","name":"symbol","type":"string"},{"indexed":false,"internalType":"uint256","name":"oldPrice","type":"uint256"},
	{"indexed":false,"internalType":"uint256","name":"newPrice","type":"uint256"}],"name":"PriceUpdated","type":"event"},
	{"inputs":[{"internalType":"string","name":"symbol","type":"string"}],"name":"getChainlinkPrice","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"internalType":"string","name":"symbol","type":"string"}],"name":"getPrice","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"internalType":"string","name":"","type":"string"}],"name":"prices","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},
	{"inputs":[{"internalType":"string","name":"symbol","type":"string"},{"internalType":"uint256","name":"newPrice","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}
]`
