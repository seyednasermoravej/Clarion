package main

import (
    "fmt"
    "math/big"
    "math/rand"
)

var prime *big.Int

func init() {
    prime, _ = big.NewInt(0).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639747", 10)
}

// GenerateCoefficients generates random coefficients for the polynomial
func GenerateCoefficients(secret *big.Int, threshold, numShares int) []*big.Int {
    coefficients := make([]*big.Int, threshold)
    coefficients[0] = secret

    for i := 1; i < threshold; i++ {
        coefficients[i] = big.NewInt(int64(rand.Intn(256)))
    }

    return coefficients
}

// EvaluatePolynomial evaluates the polynomial at the given x value
func EvaluatePolynomial(coefficients []*big.Int, x *big.Int) *big.Int {
    result := big.NewInt(0)
    temp := big.NewInt(1)

    for _, coefficient := range coefficients {
        term := big.NewInt(0).Set(coefficient)
        term.Mul(term, temp)
        result.Add(result, term)
        result.Mod(result, prime)
        temp.Mul(temp, x)
        temp.Mod(temp, prime)
    }

    return result
}

// SplitSecret splits the secret into shares
func SplitSecret(secret []byte, threshold, numShares int) []*big.Int {
    secretBigInt := new(big.Int).SetBytes(secret)
    coefficients := GenerateCoefficients(secretBigInt, threshold, numShares)

    shares := make([]*big.Int, numShares)
    for i := 1; i <= numShares; i++ {
        x := big.NewInt(int64(i))
        shares[i-1] = EvaluatePolynomial(coefficients, x)
    }

    return shares
}

// LagrangeInterpolation interpolates the secret using the given shares
func LagrangeInterpolation(shares []*big.Int, x *big.Int) *big.Int {
    result := big.NewInt(0)

    for i, share := range shares {
        numerator := big.NewInt(1)
        denominator := big.NewInt(1)

        for j, _ := range shares {
            if i != j {
                numerator.Mul(numerator, big.NewInt(0).Sub(x, big.NewInt(int64(j+1))))
                denominator.Mul(denominator, big.NewInt(0).Sub(big.NewInt(int64(i+1)), big.NewInt(int64(j+1))))
            }
        }

        lambda := big.NewInt(0).Div(numerator, denominator)
        lambda.Mod(lambda, prime)

        term := big.NewInt(0).Set(share)
        term.Mul(term, lambda)
        result.Add(result, term)
        result.Mod(result, prime)
    }

    return result
}

// ReconstructSecret reconstructs the secret from the shares
func ReconstructSecret(shares []*big.Int) string {
    secret := LagrangeInterpolation(shares, big.NewInt(0))
    secretBytes := secret.Bytes()
    return string(secretBytes)
}

func shamirTest() {
    // Define parameters
    threshold := 5           // Number of shares required to reconstruct the secret
    numShares := 9          // Total number of shares
    msg := "Hello World!" // Secret to be shared
    secret:= []byte(msg)

    // Split the secret into shares
    shares := SplitSecret(secret, threshold, numShares)
    fmt.Println("Shares:", shares)

    // Reconstruct the secret from a subset of shares
    reconstructedSecret := string(ReconstructSecret(shares[:threshold]))
    fmt.Println("Reconstructed Secret:", reconstructedSecret)
}