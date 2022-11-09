package cryptography

import (
	"fmt"
	"math/big"
	"crypto/rand"
)

type curve struct {
	p big.Int
	a big.Int
	b big.Int
}

type point struct {
	curve *curve
	x big.Int
	y big.Int
}

type generator struct {
	G point
	order big.Int
}

// Signature class for ECDSA
type Signature struct {
	r *big.Int
	s *big.Int
}

// PublicKey Class wrapping a public key (a point on the curve) and utility methods
type PublicKey struct {
	p point
}

var secp256k1 curve = curve{
	p: hexToBigInt("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F"),
	a: hexToBigInt("0x0000000000000000000000000000000000000000000000000000000000000000"),
	b: hexToBigInt("0x0000000000000000000000000000000000000000000000000000000000000007"),
}

var g point = point{
	curve: &secp256k1,
	x: hexToBigInt("0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798"),
	y: hexToBigInt("0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"),
}

var gen generator = generator{
	G: g,
	order: hexToBigInt("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141"),
}

// SignMessage using ECDSA (non-deterministic!)
func SignMessage(secretKey *big.Int, message []byte) Signature {
	e := new(big.Int).SetBytes(Hash256(message))
	n := &gen.order

	k := gen.randomSecretKey()
	curvePoint := gen.G.curve.curveMultiply(k, gen.G)

	r := new(big.Int)
	*r = curvePoint.x

	s := new(big.Int)
	s.Mul(r, secretKey).Add(s, e).Mul(s, new(big.Int).ModInverse(k, n)).Mod(s, n)

	if new(big.Int).Div(n, big.NewInt(2)).Cmp(s) == -1 {
		s.Neg(s).Add(n, s)
	}

	return Signature{r: r, s: s}
}

// VerifySignature using ECDSA
func (sig *Signature) VerifySignature(pk PublicKey, message []byte) bool {
	n := &gen.order

	if !pk.isValidPublicKey() || !sig.isValidSignature() {
		return false
	}

	e := new(big.Int).SetBytes(Hash256(message))

	sInv := new(big.Int)
	sInv.ModInverse(sig.s, n)

	u1 := new(big.Int)
	u1.Mul(sInv, e).Mod(u1, n)

	u2 := new(big.Int)
	u2.Mul(sig.r, sInv).Mod(u2, n)

	u1G := pk.p.curve.curveMultiply(u1, gen.G)
	u2Q := pk.p.curve.curveMultiply(u2, pk.p)

	curvePoint := pk.p.curve.addPointsOnCurve(u1G, u2Q)

	if curvePoint.isZero() {
		fmt.Println("Curve point is zero")
		return false
	}

	test := new(big.Int)
	test.Mod(&curvePoint.x, n)

	if test.Cmp(sig.r) == 0 {
		// r = x_1 (mod n)
		return true
	}
	
	fmt.Println("Failed equivalence test")
	return false
}

func (sig Signature) isValidSignature() bool {
	if sig.r.Cmp(big.NewInt(0)) == 1 && sig.s.Cmp(big.NewInt(0)) == 1 &&
		 sig.r.Cmp(&gen.order) == -1 && sig.s.Cmp(&gen.order) == -1 {
			return true
	}
	return false
}

func (pk PublicKey) isValidPublicKey() bool {
	// Is valid pubkey, Q!=0, Q on curve, nQ=n(cG)=0
	p := pk.p
	if p.isZero() || !p.isOnCurve() {
		return false
	}
	if !p.curve.curveMultiply(&gen.order, p).isZero() {
		return false
	}

	return true
}

func (p point) isOnCurve() bool {
	// Points on curve satisfy y^2 = x^3 + a*x + b (mod p)
	y2, x3, aX, result := new(big.Int), new(big.Int), new(big.Int), new(big.Int)

	y2.Exp(&p.y, big.NewInt(2), nil)
	x3.Exp(&p.x, big.NewInt(3), nil)
	aX.Mul(&p.x, &p.curve.a)

	result.Add(result, y2).Sub(result, x3).Sub(result, aX).Sub(result, &p.curve.b)
	result.Mod(result, &p.curve.p)

	return result.Cmp(big.NewInt(0)) == 0
}

// YfromX recovers the two possible values of the y-coord on secp256k1 from x
func YfromX(x *big.Int) (even, odd *big.Int) {
	ySq := new(big.Int)
	ySq.Exp(x, big.NewInt(3), &secp256k1.p).Add(ySq, big.NewInt(7)).Mod(ySq, &secp256k1.p)

	// Works for secp256k1 in particular
	// (https://bitcoin.stackexchange.com/questions/86234/how-to-uncompress-a-public-key)
	exp := new(big.Int)
	exp.Add(&secp256k1.p, big.NewInt(1)).Div(exp, big.NewInt(4))

	y1, y2 := new(big.Int), new(big.Int)
	y1.Exp(ySq, exp, &secp256k1.p)
	y2.Sub(&secp256k1.p, y1)

	remainder := new(big.Int).Mod(y1, big.NewInt(2))
	if remainder.Cmp(big.NewInt(0)) == 0 {
		return y1, y2
	}
	return y2, y1
}

func (p point) isZero() bool {
	zero := big.NewInt(0)
	return p.x.Cmp(zero) == 0 && p.y.Cmp(zero) == 0
}

func (gen generator) randomSecretKey() *big.Int {
	// Generate cryptographically strong pseudo-random secret key between 1 <= n <= order
	n, err := rand.Int(rand.Reader, &gen.order)
	if err != nil {
		fmt.Println(err)
	}

	return n
}

// RandomKeyPair provides a random secret key, public key pair
func RandomKeyPair() (secret *big.Int, pk PublicKey) {
	secret = gen.randomSecretKey()
	pk = gen.publicKeyFromSecretKey(secret)
	return
}

func (gen generator) publicKeyFromSecretKey(secret *big.Int) PublicKey {
	return PublicKey{p: gen.G.curve.curveMultiply(secret, gen.G)}
}

func (c *curve) curveMultiply(n *big.Int, p point) point {
	// Double and add algorithm
	bin := fmt.Sprintf("%b", n)

	result := point{curve: c, x: *big.NewInt(0), y: *big.NewInt(0)}
	append := p
	for i := len(bin) - 1; i >= 0; i-- {
		if bin[i] == '1' {
			result = c.addPointsOnCurve(result, append)
		}
		append = c.addPointsOnCurve(append, append)
	}

	return result
}

func (c *curve) addPointsOnCurve(a, b point) point {
	// Compose two points via elliptic curve addition and return result

	// A + 0 = A = 0 + A
	if a.isZero() {
		return b
	}
	if b.isZero() {
		return a
	}
	
	// A + (-A) = 0
	if a.x.Cmp(&b.x) == 0 && a.y.Cmp(&b.y) != 0 {
		return point{curve: c, x: *big.NewInt(0), y: *big.NewInt(0)}
	}

	p := &c.p

	m := new(big.Int)
	if a.x.Cmp(&b.x) == 0 { // a == b
		n := new(big.Int)
		n.Mul(big.NewInt(2), &a.y).ModInverse(n, p)

		coeff := new(big.Int)
		coeff.Exp(&a.x, big.NewInt(2), nil).Mul(coeff, big.NewInt(3)).Add(coeff, &c.a)

		m.Mul(coeff, n)
	} else {
		n := new(big.Int)
		n.Sub(&a.x, &b.x)
		n.ModInverse(n, p)

		coeff := new(big.Int)
		coeff.Sub(&a.y, &b.y)

		m.Mul(coeff, n)
	}

	newX, newY := new(big.Int), new(big.Int)
	newX.Exp(m, big.NewInt(2), nil).Sub(newX, &a.x).Sub(newX, &b.x).Mod(newX, p)
	newY.Sub(newX, &a.x).Mul(newY, m).Add(newY, &a.y).Neg(newY).Mod(newY, p)

	return point{curve: c, x: *newX, y: *newY}
}

func hexToBigInt(hex string) big.Int {
	i, _ := new(big.Int).SetString(hex, 0)
	return *i
}