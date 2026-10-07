# boss-api

The contract of BOSS's read-only JSON API, `/api/v1` on
[trumae.tech](https://trumae.tech): request paths and the Go types of
every response, in one dependency-free package. The site and its
clients, starting with [BOSSMCP](https://github.com/OptionsBoss/bossmcp),
import this package, so they can't drift apart.

```go
import bossapi "github.com/OptionsBoss/boss-api"

var s bossapi.Strategy // GET /api/v1/strategies/{slug}
```

## Conventions

- **GET only.** Nothing is ever changed, no cookie is set, no session
  exists.
- **Language** is `?lang=en|es|pt`. Empty means `en`.
- **Errors** come back as JSON `{"code": "...", "message": "..."}` with
  status 404 (`not_found`) or 400 (`invalid_argument`).
- **Caching:** content may be cached for 5 minutes. Market data and paper
  analysis are `no-store`.
- **Prices** are in USD. Each quote carries `underlying_usd`, so a client
  can convert back to a coin-denominated venue's unit (Deribit: divide by
  it). It also carries the venue's fee schedule: fee per contract =
  `min(fee_rate × underlying_usd, fee_cap × premium_usd)`.
- **Links** in content (`path` fields and markdown) are language-less site
  paths such as `/strategy/long-call`. Prefix `/es` or `/pt` for those
  languages.
- **Market data is mainnet's.** Testnets list similar instruments at
  different prices.

## Endpoints

| Path | Parameters | Response |
|---|---|---|
| `/api/v1/strategies` | `lang` | `StrategyList` |
| `/api/v1/strategies/{slug}` | `lang` | `Strategy` (with markdown) |
| `/api/v1/articles` | `lang` | `ArticleList` (published articles only) |
| `/api/v1/articles/{group}/{slug}` | `lang` | `Article` |
| `/api/v1/concepts/{key}` | `lang` | `Article` explaining the concept |
| `/api/v1/glossary/{term}` | `lang` | `GlossaryEntry` (term or concept key) |
| `/api/v1/search` | `q`, `lang` | `SearchResults` (strategies, articles, glossary) |
| `/api/v1/markets` | — | `MarketList`: venues and underlyings tracked |
| `/api/v1/expiries` | `exchange`, `currency` | `ExpiryList` |
| `/api/v1/chain` | `exchange`, `currency`, `expiry` (YYYY-MM-DD) | `Chain` |
| `/api/v1/quote` | `exchange`, `instrument` | `Quote` |
| `/api/v1/scan` | `exchange`, `currency`, `direction`, `proficiency`, `strategy`, `min_dte`, `max_dte`, `sort`, `desc`, `limit` (≤20), `lang` | `ScanResult`: the scanner's ranking after execution costs, each candidate with its legs |
| `/api/v1/paper/analyze` | `p` (a paper position, JSON) **or** `exchange`, `currency`, `legs` (JSON `[]PaperLegInput`), `qty`, `collateral`, `margin_model`; `lang` | `PaperAnalysis`: P&L, Greeks, margin, max profit/loss, breakevens, probability of profit, signals |
| `/api/v1/paper/adjust` | the same position, plus `mode` (`adjust`/`roll`), `slug`, `limit` | `PaperAdjustments`: priced alternatives |

`PaperAnalysis.position` is the position in the site's paper format. Pass
it back as `p` to analyze it again or price its adjustments.

### Opening a position in the site's paper trading

`https://trumae.tech/{lang}/paper#share=` + base64url (no padding) of
UTF-8 JSON `{"id", "slug", "name", "query", "pos"}`. `pos` is a paper
position; `query` is a strategy-page query such as
`currency=BTC&exchange=deribit`; `slug` may be empty. The page asks the
user before importing the position into their browser. Nothing is sent
to the server.

## Compatibility

Fields are only ever added. Changing the meaning of a field, or removing
one, means `/api/v2` and a new major version of this module.

## License

MIT. See [LICENSE](LICENSE).
