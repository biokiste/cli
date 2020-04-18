# migrate
script (cli) to migrate dbs

## setup
add `config.toml` with following params:

`connectionDeprecatedDB = ".../foodkoop_bk"`

`connectionDB = ".../bk_dev"`

`user_states = ["active", "former", "paused"]`

`transaction_states = ["accepted", "open", "rejected"]`

`transaction_types = ["payment", "deposit", "correction", "percent19", "percent7", "paymentSepa"]`

`token = "auth0 user token"`

`api_base_url = "http://localhost:1316/"`

`auth0URI = "https://biokiste.eu.auth0.com/"`

`audience = "https://biokiste.eu.auth0.com/api/v2/"`

`clientId = "..."`

`clientSecret = "..."`

## execute
Comment out optionally `RemoveAuthUser`, `AddUserReq` or `AddUserTransaction`
to not perform operation.
