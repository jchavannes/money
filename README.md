# Money Tracker

### Application Layers

| Layer | This project | Description | Examples |
| ----- | ------------ | ----------- | -------- |
| Template | jQuery (Javascript) | Dynamic presentation | InvestmentTransactions, Portfolio, Chart |
| Form | AJAX, WS, jQuery (Javascript) | Dynamically loading data to front-end | AddInvestmentTransaction, Signup, UpdateInvestment |
| Section | JSApp (Javascript) | Initialize sections of page | InvestmentTransactions, Portfolio, Chart |
| HTML | HTML (Go Template) | Static presentation | dashboard.html, singup.html |
| HTTP | JGo web (Go) | Handle requests and responses | /dashboard, /signup, /signup-submit |
| Domain | raw (Go) | High-level concepts | auth, chart |
| Object | raw (Go) | Collection of data | portfolio, price |
| Data | gorm (Go) | Raw data | db (investment, investment_price, investment_transaction, session, user) |
