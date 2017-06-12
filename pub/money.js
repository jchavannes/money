var MoneyApp = {};

$(function () {

    /**
     * @param token {string}
     */
    function InitCsrf(token) {
        /**
         * @param method {string}
         * @returns {boolean}
         */
        function csrfSafeMethod(method) {
            // HTTP methods that do not require CSRF protection.
            return (/^(GET|HEAD|OPTIONS|TRACE)$/.test(method));
        }

        $.ajaxSetup({
            crossDomain: false,
            beforeSend: function (xhr, settings) {
                if (!csrfSafeMethod(settings.type)) {
                    xhr.setRequestHeader("X-CSRF-Token", token);
                }
            }
        });
    }

    MoneyApp.InitCsrf = InitCsrf;

    var BaseURL = "/";

    /**
     * @param url {string}
     */
    function SetBaseUrl(url) {
        BaseURL = url;
    }

    MoneyApp.SetBaseUrl = SetBaseUrl;

    MoneyApp.Section = {
        /**
         * @param {jQuery} $portfolio
         */
        Portfolio: function ($portfolio) {
            $.ajax({
                url: MoneyApp.URL.PortfolioGet,
                method: "post",
                /**
                 * @param {string} data
                 */
                success: function (data) {
                    /** @type {Portfolio} portfolio */
                    var portfolio;
                    try {
                        portfolio = JSON.parse(data);
                    } catch (e) {
                        console.log(e);
                        return;
                    }

                    MoneyApp.Templates.Portfolio($portfolio, portfolio);
                }
            })
        },
        /**
         * @param {jQuery} $investmentTransactions
         */
        InvestmentTransactions: function ($investmentTransactions) {
            $.ajax({
                url: MoneyApp.URL.InvestmentTransactionsGet,
                method: "post",
                /**
                 * @param {string} data
                 */
                success: function (data) {
                    /** @type {[Transaction]} transactions */
                    var transactions;

                    try {
                        transactions = JSON.parse(data);
                    } catch (e) {
                        console.log(e);
                        return;
                    }
                    MoneyApp.Templates.InvestmentTransactions($investmentTransactions, transactions);
                }
            })
        }
    };

    MoneyApp.Form = {
        /**
         * @param {jQuery} $form
         * @param {int} transactionId
         */
        DeleteInvestmentTransaction: function ($form, transactionId) {
            $form.submit(function (e) {
                e.preventDefault();
                if (!confirm("Are you sure you want to delete this investment transaction?")) {
                    return;
                }
                $.ajax({
                    method: "post",
                    url: MoneyApp.URL.InvestmentTransactionDelete,
                    data: {
                        transactionId: transactionId
                    },
                    success: function () {
                        Events.Publish(MoneyApp.Events.UpdateInvestmentTransactions, {});
                    }
                });
            })
        },
        /**
         * @param {jQuery} $form
         */
        AddInvestmentTransaction: function ($form) {
            var $investmentTypeInput = $form.find("[name=investment-type]");
            var $investmentSymbolInput = $form.find("[name=investment-symbol]");
            var $transactionDateInput = $form.find("input[name=transaction-date]");
            var $transactionTypeInput = $form.find("[name=transaction-type]");
            var $transactionPriceInput = $form.find("input[name=transaction-price]");
            var $transactionQuantityInput = $form.find("input[name=transaction-quantity]");

            $transactionDateInput.datepicker();
            setInventorySymbol();
            $investmentTypeInput.change(function () {
                setInventorySymbol();
            });

            $form.submit(function (e) {
                e.preventDefault();
                $.ajax({
                    url: MoneyApp.URL.InvestmentTransactionAdd,
                    method: "post",
                    data: {
                        type: $investmentTypeInput.val(),
                        symbol: $investmentSymbolInput.val(),
                        "transaction-type": $transactionTypeInput.val(),
                        date: $transactionDateInput.val(),
                        price: $transactionPriceInput.val(),
                        quantity: $transactionQuantityInput.val()
                    },
                    success: function () {
                        Events.Publish(MoneyApp.Events.UpdateInvestmentTransactions, {});
                    }
                });
            });

            function setInventorySymbol() {
                var investmentType = $investmentTypeInput.val();
                $.ajax({
                    method: "post",
                    url: MoneyApp.URL.InvestmentSymbolsGet,
                    data: {
                        type: investmentType
                    },
                    success: function (data) {
                        var tags;
                        try {
                            /** @type {[string]} investments */
                            tags = JSON.parse(data);
                        } catch (e) {
                            console.log(e);
                            return;
                        }
                        if (tags && tags.length > 0) {
                            $investmentSymbolInput.autocomplete({
                                source: tags
                            });
                        }
                    }
                });
                var placeholder = "";
                switch (investmentType) {
                    case "crypto":
                        placeholder = "bitcoin";
                        break;
                    case "nasdaq":
                        placeholder = "AAPL";
                        break;
                    case "nyse":
                        placeholder = "WMT";
                        break;
                    case "indexsp":
                        placeholder = "INX";
                        break;
                }
                $investmentSymbolInput.attr("placeholder", placeholder);
            }
        },
        /**
         * @param {jQuery} $ele
         */
        Signup: function ($ele) {
            $ele.submit(function (e) {
                e.preventDefault();
                var username = $ele.find("[name=username]").val();
                var password = $ele.find("[name=password]").val();

                if (username.length == 0) {
                    alert("Must enter a username.");
                    return;
                }

                if (password.length == 0) {
                    alert("Must enter a password.");
                    return;
                }

                $.ajax({
                    type: "POST",
                    url: BaseURL + MoneyApp.URL.SignupSubmit,
                    data: {
                        username: username,
                        password: password
                    },
                    success: function () {
                        window.location = BaseURL + MoneyApp.URL.Dashboard
                    },
                    /**
                     * @param {XMLHttpRequest} xhr
                     */
                    error: function (xhr) {
                        alert("Error creating account:\n" + xhr.responseText + "\nIf this problem persists, try refreshing the page.");
                    }
                });
            });
        },
        /**
         * @param {jQuery} $ele
         */
        Login: function ($ele) {
            $ele.submit(function (e) {
                e.preventDefault();
                var username = $ele.find("[name=username]").val();
                var password = $ele.find("[name=password]").val();

                if (username.length == 0) {
                    alert("Must enter a username.");
                    return;
                }

                if (password.length == 0) {
                    alert("Must enter a password.");
                    return;
                }

                $.ajax({
                    type: "POST",
                    url: BaseURL + MoneyApp.URL.LoginSubmit,
                    data: {
                        username: username,
                        password: password
                    },
                    success: function () {
                        window.location = BaseURL + MoneyApp.URL.Dashboard
                    },
                    /**
                     * @param {XMLHttpRequest} xhr
                     */
                    error: function (xhr) {
                        alert("Error logging in:\n" + xhr.responseText + "\nIf this problem persists, try refreshing the page.");
                    }
                });
            });
        }
    };

    MoneyApp.Templates = {
        /**
         * @param {jQuery} $portfolio
         * @param {Portfolio} portfolio
         */
        Portfolio: function ($portfolio, portfolio) {
            var item;
            var i;
            var html = "";
            for (i = 0; i < portfolio.Items.length; i++) {
                /** @type {PortfolioItem} item */
                item = portfolio.Items[i];
                html +=
                    "<tr>" +
                    "<td>" + item.Investment.Symbol + "</td>" +
                    "<td>" + item.Quantity + "</td>" +
                    "<td>" + item.Price + "</td>" +
                    "<td>" + item.Value + "</td>" +
                    "<td>" + item.Cost + "</td>" +
                    "<td>" + item.NetGainLoss + "</td>" +
                    "<td>" + item.NetGainLossPercent + "</td>" +
                    "<td>" + item.DistributionPercent + "</td>" +
                    "<td>" + item.NetGainLossWeighted + "</td>" +
                    "</tr>";
            }
            html =
                "<table class='table table-bordered table-striped'>" +
                "<thead>" +
                "<tr>" +
                "<th>Name</th>" +
                "<th>Qty</th>" +
                "<th>Price</th>" +
                "<th>Value</th>" +
                "<th>Paid</th>" +
                "<th>Net</th>" +
                "<th>Change</th>" +
                "<th>Dist.</th>" +
                "<th>Weighted</th>" +
                "</tr>" +
                "</thead>" +
                "<tbody>" +
                html +
                "</tbody>" +
                "</table>";
            html = MoneyApp.Templates.Snippets.Panel("Portfolio", html);
            $portfolio.html(html);
        },
        /**
         * @param {jQuery} $investmentTransactions
         * @param {[Transaction]} transactions
         */
        InvestmentTransactions: function ($investmentTransactions, transactions) {
            var transaction;
            var i;
            var html = "";
            for (i = 0; i < transactions.length; i++) {
                transaction = transactions[i];
                html +=
                    "<tr>" +
                    "<td>" + (transaction.Type === 1 ? "Buy" : "Sell") + "</td>" +
                    "<td>" + transaction.Date.slice(0, 10) + "</td>" +
                    "<td>" + transaction.Investment.InvestmentType.toUpperCase() + "</td>" +
                    "<td>" + transaction.Investment.Symbol + "</td>" +
                    "<td>" + transaction.Price + "</td>" +
                    "<td>" + transaction.Quantity + "</td>" +
                    "<td>" +
                    "<form id='delete-transaction-" + transaction.Id + "'>" +
                    "<input type='submit' class='btn btn-xs btn-danger' value='Remove'/>" +
                    "</form>" +
                    "</td>" +
                    "</tr>";
            }
            html =
                "<table class='table table-bordered table-striped'>" +
                "<thead>" +
                "<tr>" +
                "<th>Type</th>" +
                "<th>Date</th>" +
                "<th>Market</th>" +
                "<th>Name</th>" +
                "<th>Price</th>" +
                "<th>Quantity</th>" +
                "<th>Actions</th>" +
                "</tr>" +
                "</thead>" +
                "<tbody>" +
                html +
                "</tbody>" +
                "</table>";
            html = MoneyApp.Templates.Snippets.Panel("Investment Transactions", html);
            $investmentTransactions.html(html);
            for (i = 0; i < transactions.length; i++) {
                transaction = transactions[i];
                MoneyApp.Form.DeleteInvestmentTransaction($("#delete-transaction-" + transaction.Id), transaction.Id);
            }
        },
        Snippets: {
            /**
             * @param {string} title
             * @param {string} html
             * @return {string}
             */
            Panel: function (title, html) {
                html =
                    "<div class='panel panel-default'>" +
                    "<div class='panel-heading'><h3 class='panel-title'>" + title + "</h3></div>" +
                    "<div class='panel-body'>" +
                    html +
                    "</div>" +
                    "</div>";
                return html;
            }
        }
    };

    MoneyApp.Events = {
        UpdateInvestmentTransactions: "update-investment-transactions"
    };

    MoneyApp.URL = {
        Dashboard: "dashboard",
        LoginSubmit: "login-submit",
        SignupSubmit: "signup-submit",
        PortfolioGet: "portfolio-get",
        InvestmentTransactionsGet: "investment-transactions-get",
        InvestmentTransactionAdd: "investment-transaction-add",
        InvestmentTransactionDelete: "investment-transaction-delete",
        InvestmentSymbolsGet: "investment-symbols-get"
    };

});

/**
 * @typedef {{
 *   Id: int
 *   Type: string
 *   Date: string
 *   Investment: Investment
 *   Quantity: float
 *   Price: float
 * }} Transaction
 */

/**
 * @typedef {{
 *   InvestmentType: string
 *   Symbol: string
 * }} Investment
 */

/**
 * @typedef {{
 *   Items: []PortfolioItem
 *   TotalValue: float
 *   TotalCost: float
 *   NetGainLoss: float
 *   NetGainLossPercent: float
 * }} Portfolio
 */

/**
 * @typedef {{
 *   Investment: Investment
 *   Url: string
 *   Quantity: float
 *   Price: float
 *   Value: float
 *   Cost: float
 *   NetGainLoss: float
 *   NetGainLossPercent: float
 *   DistributionPercent: float
 *   NetGainLossWeighted: float
 * }} PortfolioItem
 */
