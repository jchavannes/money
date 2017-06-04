var Events = {};
(function () {
    var subscriptions = {};
    Events.Subscribe = Subscribe;
    Events.Publish = Publish;
    /**
     * @param {string} eventName
     * @param {function} subscription
     */
    function Subscribe(eventName, subscription) {
        if (!subscriptions.hasOwnProperty(eventName)) {
            subscriptions[eventName] = [];
        }
        subscriptions[eventName].push(subscription);
    }

    /**
     * @param {string} eventName
     * @param {{}} data
     */
    function Publish(eventName, data) {
        if (!subscriptions.hasOwnProperty(eventName)) {
            return;
        }
        for (var i = 0; i < subscriptions[eventName].length; i++) {
            subscriptions[eventName][i](data);
        }
    }
})();
