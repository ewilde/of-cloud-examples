"use strict"
var cowsay = require("cowsay");

module.exports = (context, callback) => {
    callback(undefined, cowsay.say({
        text : context,
        e : "oO",
        T : "U "
    }));
}
