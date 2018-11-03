const QuantaForwarder = artifacts.require("QuantaForwarder");

module.exports = function(callback) {
    dest = process.argv[4].trim()
    quanta = process.argv[5].trim()

    console.log("Creating for ",dest, quanta);

    return QuantaForwarder.new(dest, quanta).then((forwarder) => {
        console.log(forwarder.address)
        callback()
    })
}