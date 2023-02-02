// BarrettMu, a class for performing Barrett modular reduction computations in
// JavaScript.
//
// Requires newBigInt.js.
//
// Copyright 2004-2005 David Shapiro.
//
// You may use, re-use, abuse, copy, and modify this code to your liking, but
// please keep this header.
//
// Thanks!
//
// Dave Shapiro
// dave@ohdave.com
const {
  biCopy,
  biHighIndex,
  biDivide,
  newBigInt,
  biMultiply,
  biDivideByRadixPower,
  biModuloByRadixPower,
  biSubtract,
  biCompare,
  biShiftRight,
  biAdd,
} = require('./BigInt.js')

function BarrettMu(m) {
  let BarrettMu = {}
  BarrettMu.modulus = biCopy(m)
  BarrettMu.k = biHighIndex(BarrettMu.modulus) + 1
  var b2k = newBigInt()
  b2k.digits[2 * BarrettMu.k] = 1 // b2k = b^(2k)
  BarrettMu.mu = biDivide(b2k, BarrettMu.modulus)
  BarrettMu.bkplus1 = newBigInt()
  BarrettMu.bkplus1.digits[BarrettMu.k + 1] = 1 // bkplus1 = b^(k+1)
  BarrettMu.modulo = BarrettMu_modulo
  BarrettMu.multiplyMod = BarrettMu_multiplyMod
  BarrettMu.powMod = BarrettMu_powMod
  return BarrettMu
}

function BarrettMu_modulo(x) {
  var q1 = biDivideByRadixPower(x, this.k - 1)
  var q2 = biMultiply(q1, this.mu)
  var q3 = biDivideByRadixPower(q2, this.k + 1)
  var r1 = biModuloByRadixPower(x, this.k + 1)
  var r2term = biMultiply(q3, this.modulus)
  var r2 = biModuloByRadixPower(r2term, this.k + 1)
  var r = biSubtract(r1, r2)
  if (r.isNeg) {
    r = biAdd(r, this.bkplus1)
  }
  var rgtem = biCompare(r, this.modulus) >= 0
  while (rgtem) {
    r = biSubtract(r, this.modulus)
    rgtem = biCompare(r, this.modulus) >= 0
  }
  return r
}

function BarrettMu_multiplyMod(x, y) {
  /*
	x = this.modulo(x);
	y = this.modulo(y);
	*/
  var xy = biMultiply(x, y)
  return this.modulo(xy)
}

function BarrettMu_powMod(x, y) {
  var result = newBigInt()
  result.digits[0] = 1
  var a = x
  var k = y
  while (true) {
    if ((k.digits[0] & 1) != 0) result = this.multiplyMod(result, a)
    k = biShiftRight(k, 1)
    if (k.digits[0] == 0 && biHighIndex(k) == 0) break
    a = this.multiplyMod(a, a)
  }
  return result
}

module.exports = { BarrettMu }
