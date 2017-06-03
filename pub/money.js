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
                    var html = MoneyApp.Templates.Snippets.Panel("Investment Transactions", data);
                    $investmentTransactions.html(html);
                }
            })
        }
    };

    MoneyApp.Form = {
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
            $investmentTypeInput.change(setInventorySymbol);

            setInventorySymbol();

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
                        console.log(tags);
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
                        success: function (data) {

                        }
                    });
                });
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

    MoneyApp.URL = {
        Dashboard: "dashboard",
        LoginSubmit: "login-submit",
        SignupSubmit: "signup-submit",
        InvestmentTransactionsGet: "investment-transactions-get",
        InvestmentTransactionAdd: "investment-transaction-add",
        InvestmentSymbolsGet: "investment-symbols-get"
    };

});
