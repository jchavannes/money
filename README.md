# Money Tracker

### Application Layers

| Layer | This project | Description |
| ----- | ------------ | ----------- |
| Template | jQuery (Javascript) | Dynamic presentation |
| Form | AJAX, WS, jQuery (Javascript) | Dynamically loading data to front-end |
| Section | JSApp (Javascript) | Initialize sections of page |
| HTML | HTML (Go Template) | Static presentation |
| HTTP | JGo web (Go) | Handle requests and responses |
| Domain | raw (Go) | High-level concepts | 
| Object | raw (Go) | Collection of data |
| Data | gorm (Go) | Raw data |

#### Item Examples

| Type | Examples |
| ---- | ----- |
| Template | InvestmentTransactions, Portfolio, Chart |
| Form | AddInvestmentTransaction, Signup, UpdateInvestment |
| Section | InvestmentTransactions, Portfolio, Chart |
| HTML | dashboard.html, singup.html |
| HTTP | /dashboard, /signup, /signup-submit |
| Domain | auth, chart |
| Object | portfolio, price |
| Data | db (investment, investment_price, investment_transaction, session, user) |
