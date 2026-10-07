// Package bossapi is the contract of BOSS's read-only JSON API, /api/v1
// on https://trumae.tech: the request paths and the JSON shapes of every
// response. It has no dependencies, so any client -- BOSSMCP first
// (github.com/OptionsBoss/bossmcp) -- and the site itself share one
// definition and can't drift apart.
//
// Every endpoint is GET-only, sets no cookies and changes nothing. The
// language of content is the ?lang= query parameter (en, es, pt; empty
// means en). Misses and bad requests answer with an Error as JSON and a
// 404 or 400 status. Content responses may be cached for 5 minutes;
// market data (/markets, /expiries, /chain, /quote, /scan) and paper
// analysis never are.
//
// Compatibility: fields are only ever added. A change to an existing
// field's meaning, or a removal, is a new API version (/api/v2) and a new
// major version of this module.
package bossapi

import (
	"encoding/json"
	"strings"
	"time"
)

// Prefix is where the API is mounted.
const Prefix = "/api/v1/"

// Langs are the languages the API answers in (?lang=); empty means "en".
var Langs = []string{"en", "es", "pt"}

// Quote is one option's live book as BOSS holds it. Every price is in USD:
// BOSS converts a venue's coin-denominated premiums when it records them
// (Deribit: coin price x underlying_price), so a client wanting the
// venue's own unit back divides by Underlying -- the exact inverse.
type Quote struct {
	Exchange   string    `json:"exchange"`
	Currency   string    `json:"currency"`
	Instrument string    `json:"instrument"`
	Right      string    `json:"right"`
	Strike     float64   `json:"strike"`
	Expiry     time.Time `json:"expiry"`
	Bid        float64   `json:"bid_usd"`
	Ask        float64   `json:"ask_usd"`
	Mid        float64   `json:"mid_usd"`
	Mark       float64   `json:"mark_usd"`
	BidSize    float64   `json:"bid_size"`
	AskSize    float64   `json:"ask_size"`
	IV         float64   `json:"iv"`
	Underlying float64   `json:"underlying_usd"`
	IndexPrice float64   `json:"index_usd"`
	UpdatedAt  time.Time `json:"updated_at"`
	// FeeRate and FeeCap are the venue's estimated taker fee schedule:
	// fee per contract = min(FeeRate x underlying_usd, FeeCap x premium_usd).
	// The schedule rather than one fee, so a client can price any limit.
	FeeRate float64 `json:"fee_rate"`
	FeeCap  float64 `json:"fee_cap"`
}

type (
	StrategyList struct {
		Strategies []StrategySummary `json:"strategies"`
	}
	ArticleList struct {
		Articles []ArticleSummary `json:"articles"`
	}
	SearchResults struct {
		Results []SearchResult `json:"results"`
	}
)

// CurrencyOf is the underlying an instrument name belongs to, the
// currency the market endpoints take:
// "BTC-27JUN25-100000-C" and "BTC_USDC-27JUN25-100000-C" are both BTC.
func CurrencyOf(instrument string) string {
	head, _, _ := strings.Cut(instrument, "-")
	coin, _, _ := strings.Cut(head, "_")
	return coin
}

// Error is a request the API can't answer. Code is stable and
// machine-readable; Message is for a human.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }

const (
	CodeNotFound        = "not_found"
	CodeInvalidArgument = "invalid_argument"
)

// Classification is a strategy's six dimensions, localized.
type Classification struct {
	Proficiency string `json:"proficiency"`
	Direction   string `json:"direction"`
	Volatility  string `json:"volatility"`
	Risk        string `json:"risk"`
	Reward      string `json:"reward"`
	Type        string `json:"type"`
}

type StrategySummary struct {
	Slug           string         `json:"slug"`
	Name           string         `json:"name"`
	Classification Classification `json:"classification"`
	Path           string         `json:"path"`
}

type Leg struct {
	Right        string `json:"right"`
	Position     string `json:"position"`
	Ratio        int    `json:"ratio"`
	StrikeOffset int    `json:"strike_offset" jsonschema:"strikes away from the at-the-money strike on the expiry's strike ladder; 0 is ATM"`
	FarExpiry    bool   `json:"far_expiry,omitempty" jsonschema:"leg belongs to a later expiry (calendars and diagonals)"`
	// The two shape knobs BOSS's strategy page exposes; telling a reader
	// which legs they move keeps a default width from being presented as
	// part of the strategy's definition.
	WidthAdjustable bool `json:"width_adjustable,omitempty"`
	RatioAdjustable bool `json:"ratio_adjustable,omitempty"`
}

type StrategyLink struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Strategy struct {
	Slug             string         `json:"slug"`
	Name             string         `json:"name"`
	Classification   Classification `json:"classification"`
	Description      string         `json:"description"`
	Setup            []string       `json:"setup"`
	Legs             []Leg          `json:"legs"`
	Context          string         `json:"context"`
	RiskProfile      string         `json:"risk_profile"`
	Pros             []string       `json:"pros"`
	Cons             []string       `json:"cons"`
	TimeEffect       string         `json:"time_effect"`
	VolatilityEffect string         `json:"volatility_effect"`
	Related          []StrategyLink `json:"related"`
	Path             string         `json:"path"`
	Markdown         string         `json:"markdown"`
}

type ArticleSummary struct {
	Key        string   `json:"key"`
	Title      string   `json:"title"`
	Definition string   `json:"definition"`
	Concepts   []string `json:"concepts"`
	Path       string   `json:"path"`
}

type Article struct {
	Key        string         `json:"key"`
	Title      string         `json:"title"`
	Definition string         `json:"definition"`
	Markdown   string         `json:"markdown"`
	Strategies []StrategyLink `json:"strategies"`
	Path       string         `json:"path"`
}

type GlossaryEntry struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
	// Concept is the learn concept key behind the term, when there is
	// one; it's also what identifies the term in a request.
	Concept string `json:"concept,omitempty"`
	// Related is the published article explaining the term in depth.
	Related  *ArticleSummary `json:"related,omitempty"`
	Markdown string          `json:"markdown"`
}

type SearchResult struct {
	Kind    string `json:"kind" jsonschema:"strategy, article or glossary"`
	Key     string `json:"key" jsonschema:"slug for a strategy, key for an article, concept key or term for a glossary entry"`
	Name    string `json:"name"`
	Snippet string `json:"snippet"`
	Path    string `json:"path,omitempty"`
}

// Market is one exchange and the currencies BOSS tracks options for there.
type Market struct {
	Exchange   string   `json:"exchange"`
	Currencies []string `json:"currencies"`
}

type MarketList struct {
	Markets []Market `json:"markets"`
}

// Expiry is one listed expiration date of a currency's options.
type Expiry struct {
	Expiry       time.Time `json:"expiry"`
	Date         string    `json:"date" jsonschema:"YYYY-MM-DD, what /api/v1/chain takes"`
	DaysToExpiry float64   `json:"days_to_expiry"`
	ForwardUSD   float64   `json:"forward_usd"`
	Strikes      int       `json:"strikes"`
}

type ExpiryList struct {
	Exchange string   `json:"exchange"`
	Currency string   `json:"currency"`
	IndexUSD float64  `json:"index_usd"`
	Expiries []Expiry `json:"expiries"`
}

// ChainRow is one strike of an expiry; Call or Put is nil where BOSS has
// no live quote for that side.
type ChainRow struct {
	Strike float64 `json:"strike"`
	Call   *Quote  `json:"call,omitempty"`
	Put    *Quote  `json:"put,omitempty"`
}

// Chain is one expiry's options, strikes ascending. Prices are in USD like
// Quote's; ForwardUSD is the underlying the venue prices this expiry off.
type Chain struct {
	Exchange   string     `json:"exchange"`
	Currency   string     `json:"currency"`
	Expiry     time.Time  `json:"expiry"`
	Date       string     `json:"date"`
	ForwardUSD float64    `json:"forward_usd"`
	IndexUSD   float64    `json:"index_usd"`
	Rows       []ChainRow `json:"rows"`
}

// PaperLegInput is one leg of a structure to analyze as if opened now:
// the instrument, the side, how many contracts per unit, and optionally
// the price actually paid or received (USD per contract) -- for a real
// position's fills instead of today's book.
type PaperLegInput struct {
	Instrument string  `json:"instrument"`
	Side       string  `json:"side" jsonschema:"buy or sell"`
	Ratio      int     `json:"ratio,omitempty" jsonschema:"contracts per unit of the structure; default 1"`
	PriceUSD   float64 `json:"price_usd,omitempty" jsonschema:"entry price per contract in USD, if already filled; default: today's execution price"`
}

type PaperGreeks struct {
	Delta float64 `json:"delta"`
	Gamma float64 `json:"gamma"`
	Vega  float64 `json:"vega"`
	Theta float64 `json:"theta"`
	Rho   float64 `json:"rho"`
}

// PaperLeg is one leg as the paper card values it now. Prices are USD
// per contract.
type PaperLeg struct {
	Instrument    string    `json:"instrument"`
	Right         string    `json:"right"`
	Side          string    `json:"side"`
	Ratio         int       `json:"ratio"`
	Strike        float64   `json:"strike,omitempty"`
	Expiry        time.Time `json:"expiry"`
	EntryPriceUSD float64   `json:"entry_price_usd"`
	ExitPriceUSD  float64   `json:"exit_price_usd,omitempty" jsonschema:"what closing it now would fetch or cost, at the book"`
	MidUSD        float64   `json:"mid_usd,omitempty"`
	IV            float64   `json:"iv,omitempty"`
	Delta         float64   `json:"delta,omitempty"`
	Settled       bool      `json:"settled,omitempty"`
	Missing       bool      `json:"missing,omitempty" jsonschema:"live, but nothing to close it into right now"`
	ThinBook      bool      `json:"thin_book,omitempty"`
}

// PaperMetrics is a structure's shape at its first expiry: what it can
// make and lose, where it breaks even and how likely it is to profit.
type PaperMetrics struct {
	MaxProfitUSD   float64     `json:"max_profit_usd" jsonschema:"the best P&L at expiry; with profit_uncapped, only the capped side's best"`
	MaxLossUSD     float64     `json:"max_loss_usd" jsonschema:"the worst P&L at expiry, negative for a loss; with loss_uncapped, only the capped side's worst (e.g. a short put at zero) -- the other side has no bound"`
	ProfitUncapped bool        `json:"profit_uncapped,omitempty"`
	LossUncapped   bool        `json:"loss_uncapped,omitempty"`
	Breakevens     []float64   `json:"breakevens"`
	PoPPct         *float64    `json:"pop_pct,omitempty" jsonschema:"probability of profit at expiry, percent"`
	Greeks         PaperGreeks `json:"greeks"`
	InitialMargin  float64     `json:"initial_margin_usd,omitempty"`
}

type PaperMargin struct {
	Model          string  `json:"model" jsonschema:"std (standard) or pm (portfolio)"`
	CollateralUSD  float64 `json:"collateral_usd"`
	EquityUSD      float64 `json:"equity_usd"`
	InitialUSD     float64 `json:"initial_margin_usd"`
	MaintenanceUSD float64 `json:"maintenance_margin_usd"`
	UtilizationPct float64 `json:"utilization_pct" jsonschema:"maintenance use that would trigger a liquidation at 100"`
}

// PaperSignal is a fact that often leads a trader to consider adjusting
// a position -- never advice. Kind is stable; Message is for a person.
type PaperSignal struct {
	Kind       string  `json:"kind" jsonschema:"tested (a short leg's delta above 0.35), margin (use above 70%), gamma (an option near expiry with high gamma), profit (over half the max profit already captured)"`
	Instrument string  `json:"instrument,omitempty"`
	Value      float64 `json:"value"`
	Message    string  `json:"message"`
}

// PaperAnalysis is the paper trading card for one position, as data.
type PaperAnalysis struct {
	// Position is the position in the site's own paper format: pass it
	// back as ?p= to /api/v1/paper/analyze or /adjust, or open it in the
	// site's paper trading page.
	Position      json.RawMessage `json:"position"`
	Exchange      string          `json:"exchange"`
	Currency      string          `json:"currency"`
	Qty           float64         `json:"qty"`
	IndexUSD      float64         `json:"index_usd"`
	EntryCostUSD  float64         `json:"entry_cost_usd" jsonschema:"paid to open (negative: credit received), fees included"`
	CloseValueUSD float64         `json:"close_value_usd" jsonschema:"what closing everything now at the book would bring (negative: cost to close)"`
	PnLUSD        float64         `json:"pnl_usd" jsonschema:"if closed now at the book, fees included"`
	MidPnLUSD     float64         `json:"mid_pnl_usd" jsonschema:"the same at mid prices"`
	ReturnPct     *float64        `json:"return_pct,omitempty"`
	Greeks        PaperGreeks     `json:"greeks"`
	Expired       bool            `json:"expired,omitempty"`
	Closed        bool            `json:"closed,omitempty"`
	Legs          []PaperLeg      `json:"legs"`
	Metrics       *PaperMetrics   `json:"metrics,omitempty" jsonschema:"the position as it stands, at its first expiry"`
	Margin        *PaperMargin    `json:"margin,omitempty" jsonschema:"only when collateral is given"`
	Signals       []PaperSignal   `json:"signals"`
}

type PaperTrade struct {
	Action    string    `json:"action" jsonschema:"close, open or keep"`
	Right     string    `json:"right"`
	Side      string    `json:"side" jsonschema:"the leg's side in the position"`
	Strike    float64   `json:"strike,omitempty"`
	Expiry    time.Time `json:"expiry"`
	Contracts float64   `json:"contracts"`
	PriceUSD  float64   `json:"price_usd"`
	FeeUSD    float64   `json:"fee_usd"`
	Name      string    `json:"name,omitempty" jsonschema:"a linear hedge's instrument"`
	Thin      bool      `json:"thin,omitempty"`
}

// PaperCandidate is one alternative to holding the position: a roll, an
// adjustment, or closing it. Cost is what the trades pay (negative:
// credit), fees included; After is the position it leaves.
type PaperCandidate struct {
	Key       string       `json:"key"`
	Kind      string       `json:"kind" jsonschema:"keep, close, roll, or an adjustment: close, out, in, wing, later, strategy, hedge"`
	Rules     []string     `json:"rules,omitempty"`
	Expiry    time.Time    `json:"expiry"`
	CostUSD   float64      `json:"cost_usd"`
	FeesUSD   float64      `json:"fees_usd"`
	Efficient bool         `json:"efficient,omitempty" jsonschema:"no other candidate is better on cost, max loss and probability of profit at once"`
	Thin      bool         `json:"thin,omitempty"`
	After     PaperMetrics `json:"after"`
	Trades    []PaperTrade `json:"trades,omitempty"`
}

type PaperAdjustments struct {
	Mode       string           `json:"mode"`
	Keep       *PaperCandidate  `json:"keep,omitempty"`
	CloseAll   *PaperCandidate  `json:"close_all,omitempty"`
	Candidates []PaperCandidate `json:"candidates"`
	Total      int              `json:"total" jsonschema:"candidates before the limit"`
	Error      string           `json:"error,omitempty"`
}

// ScanLeg is one leg of a scanner candidate, in the shape
// /api/v1/paper/analyze takes as legs.
type ScanLeg struct {
	Instrument string    `json:"instrument"`
	Side       string    `json:"side" jsonschema:"buy or sell"`
	Ratio      int       `json:"ratio"`
	Right      string    `json:"right"`
	Strike     float64   `json:"strike,omitempty"`
	Expiry     time.Time `json:"expiry"`
}

// ScanCandidate is one structure the scanner ranked: what it is, where,
// how it scores after execution costs, and how to look at it.
type ScanCandidate struct {
	Strategy     string    `json:"strategy" jsonschema:"strategy slug"`
	StrategyName string    `json:"strategy_name"`
	Exchange     string    `json:"exchange"`
	Currency     string    `json:"currency"`
	Expiry       time.Time `json:"expiry"`
	Date         string    `json:"date"`
	DTE          int       `json:"dte" jsonschema:"whole days to expiry"`
	Strike       float64   `json:"strike" jsonschema:"anchor strike every leg is placed relative to"`
	Width        int       `json:"width,omitempty"`
	WingWidth    int       `json:"wing_width,omitempty"`
	Ratio        int       `json:"ratio,omitempty"`
	Strikes      []float64 `json:"strikes"`
	Legs         []ScanLeg `json:"legs"`
	// Costs and scores, all at execution prices (ask to buy, bid to
	// sell, fees included) -- the same numbers the strategy page shows.
	NetCostUSD       float64  `json:"net_cost_usd" jsonschema:"to enter: positive is paid, negative is a credit received; fees included"`
	ExecutionCostUSD float64  `json:"execution_cost_usd" jsonschema:"how much worse than mid that is (spread + fees)"`
	PayoffDefined    bool     `json:"payoff_defined" jsonschema:"false for calendars/diagonals, whose max profit/loss and probability aren't meaningful"`
	MaxProfitUSD     float64  `json:"max_profit_usd,omitempty"`
	MaxLossUSD       float64  `json:"max_loss_usd,omitempty" jsonschema:"negative for a loss"`
	ProfitUncapped   bool     `json:"profit_uncapped,omitempty"`
	LossUncapped     bool     `json:"loss_uncapped,omitempty"`
	RewardToRisk     *float64 `json:"reward_to_risk,omitempty" jsonschema:"max profit over capital at risk; only when both are capped"`
	PoPPct           *float64 `json:"pop_pct,omitempty" jsonschema:"probability of profit at expiry, percent"`
	Lopsided         bool     `json:"lopsided,omitempty" jsonschema:"the capped loss dwarfs what it can realistically make"`
	Delta            float64  `json:"delta"`
	AverageIV        float64  `json:"average_iv"`
	Equivalents      int      `json:"equivalents,omitempty" jsonschema:"other strategies/venues with this exact payoff that cost more to execute"`
	Variants         int      `json:"variants,omitempty" jsonschema:"other strikes/widths of this strategy at this expiry that ranked lower"`
	// Path is the language-less strategy page that reproduces this exact
	// candidate (venue, expiry, strike, shape).
	Path string `json:"path"`
}

type ScanResult struct {
	Candidates []ScanCandidate `json:"candidates"`
	Total      int             `json:"total" jsonschema:"candidates matching the filters before the limit"`
	ComputedAt time.Time       `json:"computed_at"`
	// ScannerPath is the language-less scanner page with these filters.
	ScannerPath string `json:"scanner_path"`
}
