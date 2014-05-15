package polypasshash


import (
	"fmt"
  "math/rand"
	"time"
)

type shamirSecret struct {
  threshold int
	secret string
	coefficients [][]byte
}

type share struct {
	shareNumber int
	shareBytes []byte
}

func (sh *share) IsSameShare(other share) bool {
	if sh.shareNumber != other.shareNumber {
		return false
	}

	if len(sh.shareBytes) != len(other.shareBytes) {
		return false
	}

	for i:=0; i < len(other.shareBytes); i++ {
		if sh.shareBytes[i] != other.shareBytes[i] {
			return false
		}
	}

	return true
}

var random = rand.New(rand.NewSource(time.Now().Unix()))

func ShamirSecret(threshold int, secret string) shamirSecret {
  ss := shamirSecret{threshold, secret, make([][]byte, len(secret))}

	if ss.secret != "" {
		for i, c := range ss.secret {
			secretBytes := []byte{byte(c)}
			randomBytes := make([]byte, 32)
			for i:=0; i <= 31; i++{ randomBytes[i] = byte(random.Intn(256)) }
			secretBytes = append(secretBytes, randomBytes...)
			ss.coefficients[i] = secretBytes
			fmt.Println(ss.coefficients)
		}
	}

	fmt.Println("ShamirSecret initialized!!", ss.secret)
	return ss
}

func (ss *shamirSecret) ComputeShare(x int) share {
	if x <= 0 || x > 256 {
		panic("x must be between 1 and 255")
	}

	if len(ss.coefficients) == 0 {
		panic("coefficient must be initialized")
	}

	var shareBytes []byte

	for _, thisCoefficient := range ss.coefficients {
		thisShare := f(x, thisCoefficient)
		shareBytes = append(shareBytes, byte(thisShare))
	}

	return share{x, shareBytes}
}

func (ss *shamirSecret) IsValidShare(share share) bool {
	if len(ss.coefficients) == 0 {
		panic("Must initialize coefficient before validating share")
	}

	_share := ss.ComputeShare(share.shareNumber)

	if share.IsSameShare(_share) {
		return true
	} else {
		return false
	}
}

/*
func (ss *shamirSecret) RecoverSecretData(shares []share) {
	var newShares []share

	for _, share := range shares {
		new := true
		for _, ns := range newShares {
			if ns.IsSameShare(share) {
				new = false
			}
		}
		if new == true {
			newShares = append(newShares, share)
		}
	}

	shares = newShares

	if ss.threshold > len(shares) {
		panic("Threshold is smaller than the number of unique shares")
	}

	if ss.secret != "" {
		panic("Recovering secretdata when some is stored.")
	}

	var xs []int

	for _, share := range shares {
		include := false

		for _, x := range xs {
			if x == share.shareNumber { include = true}
		}

		if include == true {
			panic("Different shares with the same first byte")
		}

	  if len(share.shareBytes) != len(shares[0].shareBytes) {
			panic("Shares have different length!")
		}

		xs = append(xs, share.shareBytes)
	}

  var myCoefficients [][]byte
	var mySecretData string

	for byteToUse:=0; byteToUse<len(shares[0][1]); i++ {
		var fxs []int

		for _, share := range shares {
			fxs = append(fxs, share[1][byteToUse])
		}

		resultingPoly := fullLagrange(xs, fxs)

    match := true
		_resultingPoly := append(resultingPoly[0:ss.threshold], make([]int, len(shares)-ss.threshold))

		for i, r:= range resultingPoly {
			if _resultingPoly[i] != r { match = false}
		}

		if match == false {
			panic("Share do not match. Cannot decode")
		}

		myCoefficients = append(myCoefficients, resultingPoly)
		mySecretData += string(resultingPoly[0])
	}

	ss.coefficients = myCoefficients
	ss.secretdata = mySecretData

	fmt.Println(ss.secretdata)
}
*/

  var GF256_EXP = []int {
       0x01, 0x03, 0x05, 0x0f, 0x11, 0x33, 0x55, 0xff,
       0x1a, 0x2e, 0x72, 0x96, 0xa1, 0xf8, 0x13, 0x35,
       0x5f, 0xe1, 0x38, 0x48, 0xd8, 0x73, 0x95, 0xa4,
       0xf7, 0x02, 0x06, 0x0a, 0x1e, 0x22, 0x66, 0xaa,
       0xe5, 0x34, 0x5c, 0xe4, 0x37, 0x59, 0xeb, 0x26,
       0x6a, 0xbe, 0xd9, 0x70, 0x90, 0xab, 0xe6, 0x31,
       0x53, 0xf5, 0x04, 0x0c, 0x14, 0x3c, 0x44, 0xcc,
       0x4f, 0xd1, 0x68, 0xb8, 0xd3, 0x6e, 0xb2, 0xcd,
       0x4c, 0xd4, 0x67, 0xa9, 0xe0, 0x3b, 0x4d, 0xd7,
       0x62, 0xa6, 0xf1, 0x08, 0x18, 0x28, 0x78, 0x88,
       0x83, 0x9e, 0xb9, 0xd0, 0x6b, 0xbd, 0xdc, 0x7f,
       0x81, 0x98, 0xb3, 0xce, 0x49, 0xdb, 0x76, 0x9a,
       0xb5, 0xc4, 0x57, 0xf9, 0x10, 0x30, 0x50, 0xf0,
       0x0b, 0x1d, 0x27, 0x69, 0xbb, 0xd6, 0x61, 0xa3,
       0xfe, 0x19, 0x2b, 0x7d, 0x87, 0x92, 0xad, 0xec,
       0x2f, 0x71, 0x93, 0xae, 0xe9, 0x20, 0x60, 0xa0,
       0xfb, 0x16, 0x3a, 0x4e, 0xd2, 0x6d, 0xb7, 0xc2,
       0x5d, 0xe7, 0x32, 0x56, 0xfa, 0x15, 0x3f, 0x41,
       0xc3, 0x5e, 0xe2, 0x3d, 0x47, 0xc9, 0x40, 0xc0,
       0x5b, 0xed, 0x2c, 0x74, 0x9c, 0xbf, 0xda, 0x75,
       0x9f, 0xba, 0xd5, 0x64, 0xac, 0xef, 0x2a, 0x7e,
       0x82, 0x9d, 0xbc, 0xdf, 0x7a, 0x8e, 0x89, 0x80,
       0x9b, 0xb6, 0xc1, 0x58, 0xe8, 0x23, 0x65, 0xaf,
       0xea, 0x25, 0x6f, 0xb1, 0xc8, 0x43, 0xc5, 0x54,
       0xfc, 0x1f, 0x21, 0x63, 0xa5, 0xf4, 0x07, 0x09,
       0x1b, 0x2d, 0x77, 0x99, 0xb0, 0xcb, 0x46, 0xca,
       0x45, 0xcf, 0x4a, 0xde, 0x79, 0x8b, 0x86, 0x91,
       0xa8, 0xe3, 0x3e, 0x42, 0xc6, 0x51, 0xf3, 0x0e,
       0x12, 0x36, 0x5a, 0xee, 0x29, 0x7b, 0x8d, 0x8c,
       0x8f, 0x8a, 0x85, 0x94, 0xa7, 0xf2, 0x0d, 0x17,
       0x39, 0x4b, 0xdd, 0x7c, 0x84, 0x97, 0xa2, 0xfd,
       0x1c, 0x24, 0x6c, 0xb4, 0xc7, 0x52, 0xf6, 0x01,
       }

     var  GF256_LOG = []int {
        0x00, 0x00, 0x19, 0x01, 0x32, 0x02, 0x1a, 0xc6,
        0x4b, 0xc7, 0x1b, 0x68, 0x33, 0xee, 0xdf, 0x03,
        0x64, 0x04, 0xe0, 0x0e, 0x34, 0x8d, 0x81, 0xef,
        0x4c, 0x71, 0x08, 0xc8, 0xf8, 0x69, 0x1c, 0xc1,
        0x7d, 0xc2, 0x1d, 0xb5, 0xf9, 0xb9, 0x27, 0x6a,
        0x4d, 0xe4, 0xa6, 0x72, 0x9a, 0xc9, 0x09, 0x78,
        0x65, 0x2f, 0x8a, 0x05, 0x21, 0x0f, 0xe1, 0x24,
        0x12, 0xf0, 0x82, 0x45, 0x35, 0x93, 0xda, 0x8e,
        0x96, 0x8f, 0xdb, 0xbd, 0x36, 0xd0, 0xce, 0x94,
        0x13, 0x5c, 0xd2, 0xf1, 0x40, 0x46, 0x83, 0x38,
        0x66, 0xdd, 0xfd, 0x30, 0xbf, 0x06, 0x8b, 0x62,
        0xb3, 0x25, 0xe2, 0x98, 0x22, 0x88, 0x91, 0x10,
        0x7e, 0x6e, 0x48, 0xc3, 0xa3, 0xb6, 0x1e, 0x42,
        0x3a, 0x6b, 0x28, 0x54, 0xfa, 0x85, 0x3d, 0xba,
        0x2b, 0x79, 0x0a, 0x15, 0x9b, 0x9f, 0x5e, 0xca,
        0x4e, 0xd4, 0xac, 0xe5, 0xf3, 0x73, 0xa7, 0x57,
        0xaf, 0x58, 0xa8, 0x50, 0xf4, 0xea, 0xd6, 0x74,
        0x4f, 0xae, 0xe9, 0xd5, 0xe7, 0xe6, 0xad, 0xe8,
        0x2c, 0xd7, 0x75, 0x7a, 0xeb, 0x16, 0x0b, 0xf5,
        0x59, 0xcb, 0x5f, 0xb0, 0x9c, 0xa9, 0x51, 0xa0,
        0x7f, 0x0c, 0xf6, 0x6f, 0x17, 0xc4, 0x49, 0xec,
        0xd8, 0x43, 0x1f, 0x2d, 0xa4, 0x76, 0x7b, 0xb7,
        0xcc, 0xbb, 0x3e, 0x5a, 0xfb, 0x60, 0xb1, 0x86,
        0x3b, 0x52, 0xa1, 0x6c, 0xaa, 0x55, 0x29, 0x9d,
        0x97, 0xb2, 0x87, 0x90, 0x61, 0xbe, 0xdc, 0xfc,
        0xbc, 0x95, 0xcf, 0xcd, 0x37, 0x3f, 0x5b, 0xd1,
        0x53, 0x39, 0x84, 0x3c, 0x41, 0xa2, 0x6d, 0x47,
        0x14, 0x2a, 0x9e, 0x5d, 0x56, 0xf2, 0xd3, 0xab,
        0x44, 0x11, 0x92, 0xd9, 0x23, 0x20, 0x2e, 0x89,
        0xb4, 0x7c, 0xb8, 0x26, 0x77, 0x99, 0xe3, 0xa5,
        0x67, 0x4a, 0xed, 0xde, 0xc5, 0x31, 0xfe, 0x18,
        0x0d, 0x63, 0x8c, 0x80, 0xc0, 0xf7, 0x70, 0x07,
     }

func fullLagurange(xs []int, fxs []int) []int {
	var returnedCoefficients []int

	for i:=0; i<len(fxs); i++ {
		thisPolynomial := []int{1}

		for j:=0; j<len(fxs); j++ {
			if i == j {continue}

			denominator := gf256Sub(xs[i], xs[j])


			thisTerm := []int{gf256Div(xs[j], denominator), gf256Div(1, denominator)}

			thisPolynomial = multiplyPolynomials(thisPolynomial, thisTerm)
		}

		thisPolynomial = multiplyPolynomials(thisPolynomial, []int{fxs[i]})

		returnedCoefficients = addPolynomials(returnedCoefficients, thisPolynomial)
	}

	return returnedCoefficients
}

func f(x int, coefsBytes []byte) int {
	if x == 0 {
		panic("Invalid share index value. Cannot be 0")
  }

	accumulator, x_i := 0, 1

  for _, c := range coefsBytes {
    accumulator = gf256Add(accumulator, gf256Mul(int(c), x_i))
		x_i = gf256Mul(x_i, x)
	}

  return accumulator

}

func multiplyPolynomials(a []int, b []int) []int {
	var resultTerms, termPadding []int

	for i:=0; i<len(b); i++ {
		bterm := b[i]
		thisValue := make([]int, len(termPadding))

		copy(thisValue, termPadding);

	  for ii:=0; ii<len(a); ii++ {
		  aterm := a[ii]
			thisValue = append(thisValue, gf256Mul(aterm, bterm))
		}

		resultTerms = addPolynomials(resultTerms, thisValue)
	  termPadding = append(termPadding, 0)
	}

	return resultTerms
}


func addPolynomials(a []int, b []int) []int {


  if len(a) < len(b) {
		padding := make([]int, len(b) - len(a))
    a = append(a, padding...)
	} else if len(a) > len(b){
		padding := make([]int, len(a) - len(b))
    b = append(b, padding...)
	}

	result := make([]int, len(a))

  for pos:=0; pos<len(a); pos++{
			result[pos] = gf256Add(a[pos], b[pos])
	}


	return result
}


func gf256Add(a int, b int) int {
	return a ^ b
}

func gf256Sub(a int, b int) int {
	return gf256Add(a, b)
}

func gf256Mul(a int, b int) int {
	if a == 0 && b == 0 {
		return 0
	} else {
		return GF256_EXP[(GF256_LOG[a] + GF256_LOG[b]) % 255]
	}
}

func gf256Div(a int, b int) int {
	if a == 0 {
		return 0
	} else if b == 0 {
		panic("ZeroDivisionError")
	} else {
		mod := (GF256_LOG[a] - GF256_LOG[b]) % 255

    // Unlike Python or Ruby, Go follows truncated division for the division of negative number
		// http://golang.org/ref/spec#Arithmetic_operators
		if mod < 0 {
			mod = mod + 255
		}
		return GF256_EXP[mod]
	}
}
