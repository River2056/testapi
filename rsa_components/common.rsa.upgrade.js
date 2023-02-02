
const { RSAKeyPair, encryptedString }  = require("./RSA.js");
const { setMaxDigits } = require("./BigInt.js");

setMaxDigits(130);//
var rsa_public_key ="b66b6d86bea95ae0fde3983e69d075ec3a8c59e2c6282c15de23125ce58841ca7a99f35a0dd3feae41c28a488da38500a223db9975aa3c1d5c93cef355597cd4b668dc8bd89fa0d26d100fffa9095528f5c430ab7ae8307de87f3c786799b1c446967158d7d03afe4e06ad9a175a30bd350a5e08323dbe827264f6066302109f";
var rsa_public_exponent ="10001" ;
var rsa_deal_key = RSAKeyPair(rsa_public_exponent,"",rsa_public_key);
/**
 * 加密处理
 * @param data
 * @returns {*}
 */

function commonEncryptedString(data) {
    var value = encodeURIComponent(data);//uri编码，支持中文(特殊字等)
    var result = encryptedString(rsa_deal_key, value);
    return result;
}

module.exports = {commonEncryptedString};
