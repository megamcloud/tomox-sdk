package interfaces

import (
	"context"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/relayer"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/ws"
)

type OrderDao interface {
	GetCollection() *mgo.Collection
	Create(o *types.Order) error
	Watch() (*mgo.ChangeStream, *mgo.Session, error)
	Update(id bson.ObjectId, o *types.Order) error
	Upsert(id bson.ObjectId, o *types.Order) error
	Delete(orders ...*types.Order) error
	DeleteByHashes(hashes ...common.Hash) error
	UpdateAllByHash(h common.Hash, o *types.Order) error
	UpdateByHash(h common.Hash, o *types.Order) error
	UpsertByHash(h common.Hash, o *types.Order) error
	GetOrderCountByUserAddress(addr common.Address) (int, error)
	GetByID(id bson.ObjectId) (*types.Order, error)
	GetByHash(h common.Hash) (*types.Order, error)
	GetByHashes(hashes []common.Hash) ([]*types.Order, error)
	GetByUserAddress(addr, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error)
	GetOpenOrdersByUserAddress(addr common.Address) ([]*types.Order, error)
	GetCurrentByUserAddress(a common.Address, limit ...int) ([]*types.Order, error)
	GetHistoryByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error)
	UpdateOrderFilledAmount(h common.Hash, value *big.Int) error
	UpdateOrderFilledAmounts(h []common.Hash, values []*big.Int) ([]*types.Order, error)
	UpdateOrderStatusesByHashes(status string, hashes ...common.Hash) ([]*types.Order, error)
	GetUserLockedBalance(account common.Address, token common.Address, p []*types.Pair) (*big.Int, error)
	UpdateOrderStatus(h common.Hash, status string) error
	GetRawOrderBook(*types.Pair) ([]*types.Order, error)
	GetOrderBook(*types.Pair) ([]map[string]string, []map[string]string, error)
	GetOrderBookInDb(*types.Pair) ([]map[string]string, []map[string]string, error)
	GetSideOrderBook(p *types.Pair, side string, sort int, limit ...int) ([]map[string]string, error)
	GetOrderBookPricePoint(p *types.Pair, pp *big.Int, side string) (*big.Int, error)
	FindAndModify(h common.Hash, o *types.Order) (*types.Order, error)
	Drop() error
	Aggregate(q []bson.M) ([]*types.OrderData, error)
	AddNewOrder(o *types.Order, topic string) error
	CancelOrder(o *types.Order, topic string) error
	GetOrders(orderSpec types.OrderSpec, sort []string, offset int, size int) (*types.OrderRes, error)
	GetOrderNonce(addr common.Address) (interface{}, error)
	GetOpenOrders() ([]*types.Order, error)
}

type StopOrderDao interface {
	Create(so *types.StopOrder) error
	Update(id bson.ObjectId, so *types.StopOrder) error
	UpdateByHash(h common.Hash, so *types.StopOrder) error
	Upsert(id bson.ObjectId, so *types.StopOrder) error
	UpsertByHash(h common.Hash, so *types.StopOrder) error
	UpdateAllByHash(h common.Hash, so *types.StopOrder) error
	GetByHash(h common.Hash) (*types.StopOrder, error)
	FindAndModify(h common.Hash, so *types.StopOrder) (*types.StopOrder, error)
	GetTriggeredStopOrders(baseToken, quoteToken common.Address, lastPrice *big.Int) ([]*types.StopOrder, error)
	Drop() error
}

type AccountDao interface {
	Create(account *types.Account) (err error)
	GetAll() (res []types.Account, err error)
	GetByID(id bson.ObjectId) (*types.Account, error)
	GetByAddress(owner common.Address) (response *types.Account, err error)
	GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error)
	GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error)
	UpdateTokenBalance(owner common.Address, token common.Address, tokenBalance *types.TokenBalance) (err error)
	UpdateBalance(owner common.Address, token common.Address, balance *big.Int) (err error)
	FindOrCreate(addr common.Address) (*types.Account, error)
	Transfer(token common.Address, fromAddress common.Address, toAddress common.Address, amount *big.Int) error
	Drop()
	GetFavoriteTokens(owner common.Address) (map[common.Address]bool, error)
	AddFavoriteToken(owner, token common.Address) error
	DeleteFavoriteToken(owner, token common.Address) error
}

type RelayerDao interface {
	Create(relayer *types.Relayer) (err error)
	GetAll() (res []types.Relayer, err error)
	GetByHost(host string) (relayer *types.Relayer, err error)
	GetByAddress(addr common.Address) (relayer *types.Relayer, err error)
	DeleteByAddress(addr common.Address) error
	UpdateByAddress(addr common.Address, relayer *types.Relayer) error
	UpdateNameByAddress(addr common.Address, name string, url string) error
}

type ConfigDao interface {
	GetSchemaVersion() uint64
	GetAddressIndex(chain types.Chain) (uint64, error)
	IncrementAddressIndex(chain types.Chain) error
	ResetBlockCounters() error
	GetBlockToProcess(chain types.Chain) (uint64, error)
	SaveLastProcessedBlock(chain types.Chain, block uint64) error
	Drop()
}

type AssociationDao interface {
	GetAssociationByChainAddress(chain types.Chain, address common.Address) (*types.AddressAssociationRecord, error)
	GetAssociationByChainAssociatedAddress(chain types.Chain, associatedAddress common.Address) (*types.AddressAssociationRecord, error)

	// save mean if there is no item then insert, otherwise update
	SaveAssociation(record *types.AddressAssociationRecord) error
	SaveDepositTransaction(chain types.Chain, sourceAccount common.Address, txEnvelope string) error
	SaveAssociationStatus(chain types.Chain, sourceAccount common.Address, status string) error
}

type WalletDao interface {
	Create(wallet *types.Wallet) error
	GetAll() ([]types.Wallet, error)
	GetByID(id bson.ObjectId) (*types.Wallet, error)
	GetByAddress(addr common.Address) (*types.Wallet, error)
	GetDefaultAdminWallet() (*types.Wallet, error)
	GetOperatorWallets() ([]*types.Wallet, error)
}

type PairDao interface {
	Create(o *types.Pair) error
	GetAll() ([]types.Pair, error)
	GetAllByCoinbase(addr common.Address) ([]types.Pair, error)
	GetActivePairs() ([]*types.Pair, error)
	GetActivePairsByCoinbase(addr common.Address) ([]*types.Pair, error)
	GetByID(id bson.ObjectId) (*types.Pair, error)
	GetByName(name string) (*types.Pair, error)
	GetByTokenSymbols(baseTokenSymbol, quoteTokenSymbol string) (*types.Pair, error)
	GetByTokenAddress(baseToken, quoteToken common.Address) (*types.Pair, error)
	GetListedPairs() ([]types.Pair, error)
	GetUnlistedPairs() ([]types.Pair, error)
	DeleteByToken(baseAddress common.Address, quoteAddress common.Address) error
	DeleteByTokenAndCoinbase(baseAddress common.Address, quoteAddress common.Address, addr common.Address) error
}

type TradeDao interface {
	GetCollection() *mgo.Collection
	Create(o ...*types.Trade) error
	Watch() (*mgo.ChangeStream, *mgo.Session, error)
	Update(t *types.Trade) error
	UpdateByHash(h common.Hash, t *types.Trade) error
	GetAll() ([]types.Trade, error)
	Aggregate(q []bson.M) ([]*types.Tick, error)
	GetByPairName(name string) ([]*types.Trade, error)
	GetByHash(h common.Hash) (*types.Trade, error)
	GetByMakerOrderHash(h common.Hash) ([]*types.Trade, error)
	GetByTakerOrderHash(h common.Hash) ([]*types.Trade, error)
	GetByOrderHashes(hashes []common.Hash) ([]*types.Trade, error)
	GetSortedTrades(bt, qt common.Address, from, to int64, n int) ([]*types.Trade, error)
	GetSortedTradesByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Trade, error)
	GetNTradesByPairAddress(bt, qt common.Address, n int) ([]*types.Trade, error)
	GetTradesByPairAddress(bt, qt common.Address, n int) ([]*types.Trade, error)
	GetAllTradesByPairAddress(bt, qt common.Address) ([]*types.Trade, error)
	FindAndModify(h common.Hash, t *types.Trade) (*types.Trade, error)
	GetByUserAddress(a common.Address) ([]*types.Trade, error)
	GetLatestTrade(bt, qt common.Address) (*types.Trade, error)
	UpdateTradeStatus(h common.Hash, status string) error
	UpdateTradeStatuses(status string, hashes ...common.Hash) ([]*types.Trade, error)
	UpdateTradeStatusesByOrderHashes(status string, hashes ...common.Hash) ([]*types.Trade, error)
	Drop()
	GetTrades(tradeSpec *types.TradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.TradeRes, error)
	GetTradesUserHistory(a common.Address, tradeSpec *types.TradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.TradeRes, error)
	GetTradeByTime(dateFrom, dateTo int64, pageOffset int, pageSize int) ([]*types.Trade, error)
}

type TokenDao interface {
	Create(token *types.Token) error
	UpdateByToken(contractAddress common.Address, token *types.Token) error
	UpdateByTokenAndCoinbase(contractAddress common.Address, addr common.Address, token *types.Token) error
	GetAll() ([]types.Token, error)
	GetAllByCoinbase(addr common.Address) ([]types.Token, error)
	GetByID(id bson.ObjectId) (*types.Token, error)
	GetByAddress(addr common.Address) (*types.Token, error)
	GetQuoteTokens() ([]types.Token, error)
	GetBaseTokens() ([]types.Token, error)
	UpdateFiatPriceBySymbol(symbol string, price float64) error
	Drop() error
	DeleteByToken(contractAddress common.Address) error
	DeleteByTokenAndCoinbase(contractAddress common.Address, addr common.Address) error
}

type NotificationDao interface {
	Create(notifications ...*types.Notification) ([]*types.Notification, error)
	GetAll() ([]types.Notification, error)
	GetByUserAddress(addr common.Address, limit int, offset int) ([]*types.Notification, error)
	GetSortDecByUserAddress(addr common.Address, limit int, offset int) ([]*types.Notification, error)
	GetByID(id bson.ObjectId) (*types.Notification, error)
	FindAndModify(id bson.ObjectId, n *types.Notification) (*types.Notification, error)
	Update(n *types.Notification) error
	Upsert(id bson.ObjectId, n *types.Notification) error
	Delete(notifications ...*types.Notification) error
	DeleteByIds(ids ...bson.ObjectId) error
	Aggregate(q []bson.M) ([]*types.Notification, error)
	Drop()
	MarkRead(id bson.ObjectId) error
	MarkUnRead(id bson.ObjectId) error
	MarkAllRead(addr common.Address) error
}

type Engine interface {
	HandleOrders(msg *rabbitmq.Message) error
	// RecoverOrders(matches types.Matches) error
	// CancelOrder(order *types.Order) (*types.EngineResponse, error)
	// DeleteOrder(o *types.Order) error
	Provider() EthereumProvider
}

type WalletService interface {
	CreateAdminWallet(a common.Address) (*types.Wallet, error)
	GetDefaultAdminWallet() (*types.Wallet, error)
	GetOperatorWallets() ([]*types.Wallet, error)
	GetOperatorAddresses() ([]common.Address, error)
	GetAll() ([]types.Wallet, error)
	GetByAddress(addr common.Address) (*types.Wallet, error)
}

type OHLCVService interface {
	Unsubscribe(c *ws.Client)
	UnsubscribeChannel(c *ws.Client, p *types.SubscriptionPayload)
	Subscribe(c *ws.Client, p *types.SubscriptionPayload)
	GetOHLCV(p []types.PairAddresses, duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error)
	Get24hTick(baseToken, quoteToken common.Address) *types.Tick
	GetFiatPriceChart() (map[string][]*types.FiatPriceItem, error)
	GetLastPriceCurrentByTime(symbol string, createAt time.Time) (*big.Float, error)
	GetAllTokenPairData() ([]*types.PairData, error)
	GetAllTokenPairDataByCoinbase(addr common.Address) ([]*types.PairData, error)
	GetTokenPairData(baseToken common.Address, quoteToken common.Address) *types.PairData
}

type EthereumService interface {
	WaitMined(hash common.Hash) (*eth.Receipt, error)
	GetPendingNonceAt(a common.Address) (uint64, error)
	GetBalanceAt(a common.Address) (*big.Int, error)
}

type OrderService interface {
	GetOrdersLockedBalanceByUserAddress(addr common.Address) (map[string]*big.Int, error)
	GetOrderCountByUserAddress(addr common.Address) (int, error)
	GetByID(id bson.ObjectId) (*types.Order, error)
	GetByHash(h common.Hash) (*types.Order, error)
	GetByHashes(hashes []common.Hash) ([]*types.Order, error)
	// GetTokenByAddress(a common.Address) (*types.Token, error)
	GetByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error)
	GetCurrentByUserAddress(a common.Address, limit ...int) ([]*types.Order, error)
	GetHistoryByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error)
	NewOrder(o *types.Order) error
	CancelOrder(oc *types.OrderCancel) error
	CancelAllOrder(a common.Address) error
	HandleEngineResponse(res *types.EngineResponse) error
	GetOrders(orderSpec types.OrderSpec, sort []string, offset int, size int) (*types.OrderRes, error)
	GetOrderNonceByUserAddress(addr common.Address) (interface{}, error)
}

type OrderBookService interface {
	GetOrderBook(bt, qt common.Address) (*types.OrderBook, error)
	GetDbOrderBook(bt, qt common.Address) (*types.OrderBook, error)
	GetRawOrderBook(bt, qt common.Address) (*types.RawOrderBook, error)
	SubscribeOrderBook(c *ws.Client, bt, qt common.Address)
	UnsubscribeOrderBook(c *ws.Client)
	UnsubscribeOrderBookChannel(c *ws.Client, bt, qt common.Address)
	SubscribeRawOrderBook(c *ws.Client, bt, qt common.Address)
	UnsubscribeRawOrderBook(c *ws.Client)
	UnsubscribeRawOrderBookChannel(c *ws.Client, bt, qt common.Address)
}

type PairService interface {
	Create(pair *types.Pair) error
	CreatePairs(token common.Address) ([]*types.Pair, error)
	GetByID(id bson.ObjectId) (*types.Pair, error)
	GetByTokenAddress(bt, qt common.Address) (*types.Pair, error)
	GetTokenPairData(bt, qt common.Address) (*types.PairData, error)
	GetAllTokenPairData() ([]*types.PairData, error)
	GetMarketStats(bt, qt common.Address) (*types.PairData, error)
	GetAllMarketStats() ([]*types.PairData, error)
	GetAll() ([]types.Pair, error)
	GetAllByCoinbase(addr common.Address) ([]types.Pair, error)
	GetListedPairs() ([]types.Pair, error)
	GetUnlistedPairs() ([]types.Pair, error)
}

type TokenService interface {
	Create(token *types.Token) error
	GetByID(id bson.ObjectId) (*types.Token, error)
	GetByAddress(a common.Address) (*types.Token, error)
	GetAll() ([]types.Token, error)
	GetAllByCoinbase(addr common.Address) ([]types.Token, error)
	GetQuoteTokens() ([]types.Token, error)
	GetBaseTokens() ([]types.Token, error)
}

type TradeService interface {
	GetByPairName(p string) ([]*types.Trade, error)
	GetAllTradesByPairAddress(bt, qt common.Address) ([]*types.Trade, error)
	GetSortedTrades(bt, qt common.Address, from, to int64, n int) ([]*types.Trade, error)
	GetSortedTradesByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Trade, error)
	GetByUserAddress(a common.Address) ([]*types.Trade, error)
	GetByHash(h common.Hash) (*types.Trade, error)
	GetByOrderHashes(h []common.Hash) ([]*types.Trade, error)
	GetByMakerOrderHash(h common.Hash) ([]*types.Trade, error)
	GetByTakerOrderHash(h common.Hash) ([]*types.Trade, error)
	Subscribe(c *ws.Client, bt, qt common.Address)
	UnsubscribeChannel(c *ws.Client, bt, qt common.Address)
	Unsubscribe(c *ws.Client)
	GetTrades(tradeSpec *types.TradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.TradeRes, error)
	GetTradesUserHistory(a common.Address, tradeSpec *types.TradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.TradeRes, error)
}

type PriceBoardService interface {
	Subscribe(c *ws.Client, bt, qt common.Address)
	UnsubscribeChannel(c *ws.Client, bt, qt common.Address)
	Unsubscribe(c *ws.Client)
}

type MarketsService interface {
	Subscribe(c *ws.Client)
	UnsubscribeChannel(c *ws.Client)
	Unsubscribe(c *ws.Client)
}

type NotificationService interface {
	Create(n *types.Notification) ([]*types.Notification, error)
	GetAll() ([]types.Notification, error)
	GetByUserAddress(a common.Address, limit int, offset int) ([]*types.Notification, error)
	GetSortDecByUserAddress(addr common.Address, limit int, offset int) ([]*types.Notification, error)
	GetByID(id bson.ObjectId) (*types.Notification, error)
	Update(n *types.Notification) (*types.Notification, error)
	MarkRead(id bson.ObjectId) error
	MarkUnRead(id bson.ObjectId) error
	MarkAllRead(addr common.Address) error
}

type TxService interface {
	GetTxCallOptions() *bind.CallOpts
	GetTxSendOptions() (*bind.TransactOpts, error)
	GetTxDefaultSendOptions() (*bind.TransactOpts, error)
	SetTxSender(w *types.Wallet)
	GetCustomTxSendOptions(w *types.Wallet) *bind.TransactOpts
}

type AccountService interface {
	GetAll() ([]types.Account, error)
	Create(account *types.Account) error
	GetByID(id bson.ObjectId) (*types.Account, error)
	GetByAddress(a common.Address) (*types.Account, error)
	FindOrCreate(a common.Address) (*types.Account, error)
	GetTokenBalance(owner common.Address, token common.Address) (*types.TokenBalance, error)
	GetTokenBalances(owner common.Address) (map[common.Address]*types.TokenBalance, error)
	Transfer(token common.Address, fromAddress common.Address, toAddress common.Address, amount *big.Int) error
	GetFavoriteTokens(account common.Address) (map[common.Address]bool, error)
	AddFavoriteToken(account, token common.Address) error
	DeleteFavoriteToken(account, token common.Address) error
	GetTokenBalanceProvidor(owner common.Address, token common.Address) (*types.TokenBalance, error)
}

type ValidatorService interface {
	ValidateBalance(o *types.Order) error
	ValidateAvailableBalance(o *types.Order) error
}

type EthereumConfig interface {
	GetURL() string
	ExchangeAddress() common.Address
}

type EthereumClient interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*eth.Receipt, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
	SendTransaction(ctx context.Context, tx *eth.Transaction) error
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	BalanceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (*big.Int, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]eth.Log, error)
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- eth.Log) (ethereum.Subscription, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

type EthereumProvider interface {
	WaitMined(h common.Hash) (*eth.Receipt, error)
	GetBalanceAt(a common.Address) (*big.Int, error)
	GetPendingNonceAt(a common.Address) (uint64, error)
	BalanceOf(owner common.Address, token common.Address) (*big.Int, error)
	Decimals(token common.Address) (uint8, error)
	Symbol(token common.Address) (string, error)
	Balance(owner common.Address, token common.Address) (*big.Int, error)
}

// RelayerService interface for relayer
type RelayerService interface {
	UpdateRelayer(addr common.Address) error
	UpdateRelayers() error
	UpdateNameByAddress(addr common.Address, name string, url string) error
	GetRelayerAddress(r *http.Request) common.Address
	GetByAddress(addr common.Address) (*types.Relayer, error)
}

// Relayer interface for relayer
type Relayer interface {
	GetRelayer(addr common.Address) (*relayer.RInfo, error)
	GetLending() (*relayer.LendingRInfo, error)
	GetRelayers() ([]*relayer.RInfo, error)
	GetLendings() ([]*relayer.LendingRInfo, error)
}

// LendingOrderService for lending
type LendingOrderService interface {
	NewLendingOrder(o *types.LendingOrder) error
	CancelLendingOrder(oc *types.LendingOrder) error
	GetLendingNonceByUserAddress(addr common.Address) (uint64, error)
	GetByHash(h common.Hash) (*types.LendingOrder, error)
	RepayLendingOrder(o *types.LendingOrder) error
	TopupLendingOrder(o *types.LendingOrder) error
	GetLendingOrders(lendingSpec types.LendingSpec, sort []string, offset int, size int) (*types.LendingRes, error)
	GetTopup(topupSpec types.TopupSpec, sort []string, offset int, size int) (*types.LendingRes, error)
	GetRepay(repaySpec types.RepaySpec, sort []string, offset int, size int) (*types.LendingRes, error)
	GetRecall(recall types.RecallSpec, sort []string, offset int, size int) (*types.LendingRes, error)
	EstimateCollateral(collateralToken common.Address, lendingToken common.Address, lendingAmount *big.Float) (*big.Float, *big.Float, error)
}

// LendingOrderDao dao
type LendingOrderDao interface {
	GetByHash(h common.Hash) (*types.LendingOrder, error)
	Watch() (*mgo.ChangeStream, *mgo.Session, error)
	GetLendingNonce(addr common.Address) (uint64, error)
	AddNewLendingOrder(o *types.LendingOrder) error
	CancelLendingOrder(o *types.LendingOrder) error
	GetLendingOrderBook(term uint64, lendingToken common.Address) ([]map[string]string, []map[string]string, error)
	GetLendingOrderBookInDb(term uint64, lendingToken common.Address) ([]map[string]string, []map[string]string, error)
	GetLendingOrderBookInterest(term uint64, lendingToken common.Address, interest uint64, side string) (*big.Int, error)
	RepayLendingOrder(o *types.LendingOrder) error
	TopupLendingOrder(o *types.LendingOrder) error
	GetLendingOrders(lendingSpec types.LendingSpec, sort []string, offset int, size int) (*types.LendingRes, error)
	GetLastTokenPrice(bToken common.Address, qToken common.Address) (*big.Int, error)
}

// LendingOrderBookService interface for lending order book
type LendingOrderBookService interface {
	GetLendingOrderBook(term uint64, lendingToken common.Address) (*types.LendingOrderBook, error)
	GetLendingOrderBookInDb(term uint64, lendingToken common.Address) (*types.LendingOrderBook, error)
	SubscribeLendingOrderBook(c *ws.Client, term uint64, lendingToken common.Address)
	UnsubscribeLendingOrderBook(c *ws.Client)
	UnsubscribeLendingOrderBookChannel(c *ws.Client, term uint64, lendingToken common.Address)
}

// LendingTradeService interface for lending service
type LendingTradeService interface {
	Subscribe(c *ws.Client, term uint64, lendingToken common.Address)
	UnsubscribeChannel(c *ws.Client, term uint64, lendingToken common.Address)
	Unsubscribe(c *ws.Client)
	GetLendingTradesUserHistory(a common.Address, lendingtradeSpec *types.LendingTradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.LendingTradeRes, error)
	GetLendingTrades(lendingtradeSpec *types.LendingTradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.LendingTradeRes, error)
	RegisterNotify(fn func(*types.LendingTrade))
	GetLendingTradeByTime(dateFrom, dateTo int64, pageOffset int, pageSize int) ([]*types.LendingTrade, error)
}

// LendingTradeDao interface for lending dao
type LendingTradeDao interface {
	GetLendingTradeByOrderBook(tern uint64, lendingToken common.Address, from, to int64, n int) ([]*types.LendingTrade, error)
	Watch() (*mgo.ChangeStream, *mgo.Session, error)
	GetLendingTradeByTime(dateFrom, dateTo int64, pageOffset int, pageSize int) ([]*types.LendingTrade, error)
	GetLendingTradesUserHistory(a common.Address, lendingtradeSpec *types.LendingTradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.LendingTradeRes, error)
	GetLendingTrades(lendingtradeSpec *types.LendingTradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.LendingTradeRes, error)
	GetByHash(hash common.Hash) (*types.LendingTrade, error)
}

// LendingOhlcvService interface for lending service
type LendingOhlcvService interface {
	GetOHLCV(term uint64, lendingToken common.Address, duration int64, unit string, timeInterval ...int64) ([]*types.LendingTick, error)
	Subscribe(conn *ws.Client, p *types.SubscriptionPayload)
	Unsubscribe(conn *ws.Client)
	GetAllTokenPairData() ([]*types.LendingTick, error)
	GetTokenPairData(term uint64, lendingToken common.Address) *types.LendingTick
}

// LendingPairDao interface for lending pair by term/lendingtoken
type LendingPairDao interface {
	Create(o *types.LendingPair) error
	GetAll() ([]types.LendingPair, error)
	GetAllByCoinbase(addr common.Address) ([]types.LendingPair, error)
	DeleteByLendingKey(term uint64, lendingAddress common.Address) error
	DeleteByLendingKeyAndCoinbase(term uint64, lendingAddress common.Address, addr common.Address) error
	GetByLendingID(term uint64, lendingToken common.Address) (*types.LendingPair, error)
}

//LendingPairService imp lending
type LendingPairService interface {
	GetAll() ([]types.LendingPair, error)
	GetAllByCoinbase(addr common.Address) ([]types.LendingPair, error)
	GetByLendingID(term uint64, lendingAddress common.Address) (*types.LendingPair, error)
}

// LendingMarketsService lending service interface
type LendingMarketsService interface {
	Subscribe(c *ws.Client)
	UnsubscribeChannel(c *ws.Client)
	Unsubscribe(c *ws.Client)
}

// LendingPriceBoardService lending price board service
type LendingPriceBoardService interface {
	Subscribe(c *ws.Client, term uint64, lendingToken common.Address)
	UnsubscribeChannel(c *ws.Client, term uint64, lendingToken common.Address)
	Unsubscribe(c *ws.Client)
}
