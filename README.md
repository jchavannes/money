# Money Tracker

### Application Layers

| Layer | This project | Description |
| ----- | ------------ | ----------- |
| Template | jQuery, Javascript | Dynamic presentation |
| Form | AJAX, Javascript | Dynamically loading data to front-end |
| Section | JSApp, Javascript | Initialize sections of page |
| HTML | Go Templates | Static presentation |
| HTTP | JGo web | Handle requests and responses |
| Domain | Chart -> Price? | High-level concepts | 
| Object | db directories (e.g. Investment) | Collection of data |
| Data | db (gorm) | Raw data |

#### Current

| Type | Items |
| ---- | ----- |
| ~~Domain~~ | _none_ |
| Object | auth, chart, investment, portfolio, price |
| Data | investment, investment_price, investment_transaction, session, user |

#### Ideal

| Type | Items |
| ---- | ----- |
| Domain | auth, chart |
| Object | portfolio, price |
| Data | investment, investment_price, investment_transaction, session, user |
