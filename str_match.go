package main

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"

	"github.com/bbalet/stopwords"
	"github.com/reiver/go-porterstemmer"
	"github.com/seedco/megophone"
)

// SimilarityScore similarity score.
type SimilarityScore struct {
	A                       string
	B                       string
	CosineSimilarityScore   float64
	AbsoluteSimilarityScore float64
	AMasterString           bool
}

// Similarity similarity finder.
func Similarity(master string, other string, useMasterStringAsCorpus bool) SimilarityScore {
	s := applyStopwords(master)
	s1 := applyStopwords(other)

	a := tokenize(s)
	b := tokenize(s1)

	corp := []string{}
	if useMasterStringAsCorpus {
		corp = corpus(a, []string{}) // corpus is formed only with the master string.
	} else {
		corp = corpus(a, b)
	}
	_, encodedAAll := encodedTokens(a)
	_, encodedBAll := encodedTokens(b)

	wordFreqA := wordFreq(encodedAAll)
	wordFreqB := wordFreq(encodedBAll)

	vectorA := toVector(corp, wordFreqA)
	vectorB := toVector(corp, wordFreqB)

	cossim, _ := cosineSimilarity(vectorA, vectorB)
	abssim := absoluteSimilarity(vectorA, vectorB)

	return SimilarityScore{A: master, B: other, CosineSimilarityScore: cossim, AbsoluteSimilarityScore: abssim, AMasterString: useMasterStringAsCorpus}

}

func main() {
	a := "Hello Kitty Kajue"
	b := "Kitty Hello Kajue Ksllow kkjusubeh shjhsuygd sjhsueb shjasjagd"

	fmt.Printf("%+v", Similarity(a, b, true))

}

func toVector(corpus []string, wordFreq map[string]int) []float64 {
	v := []float64{}

	for _, w := range corpus {
		if val, ok := wordFreq[w]; ok {
			v = append(v, float64(val))
		} else {
			v = append(v, 0)
		}
	}
	return v
}

func applyStopwords(s string) string {
	return stopwords.CleanString(s, "en", false)
}

func tokenize(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}

func corpus(a []string, b []string) []string {
	corpus := append(a, b...)

	set := make(map[string]bool)
	uniquecorpus := []string{}
	for _, t := range corpus {
		if !set[t] {
			uniquecorpus = append(uniquecorpus, strings.ToLower(t))
			x, _ := megophone.Metaphone(t)
			uniquecorpus = append(uniquecorpus, strings.ToLower(x))
			stem := porterstemmer.Stem([]rune(t))
			uniquecorpus = append(uniquecorpus, strings.ToLower(string(stem)))
			set[t] = true
		}
	}
	sort.Strings(uniquecorpus)
	return uniquecorpus
}

func encodedTokens(a []string) ([]string, []string) {
	corpus := append(a)

	set := make(map[string]bool)
	uniquecorpus := []string{}
	nonuniquecorpus := []string{}
	for _, t := range corpus {
		if !set[t] {
			uniquecorpus = append(uniquecorpus, strings.ToLower(t))
			x, _ := megophone.Metaphone(t)
			uniquecorpus = append(uniquecorpus, strings.ToLower(x))
			stem := porterstemmer.Stem([]rune(t))
			uniquecorpus = append(uniquecorpus, strings.ToLower(string(stem)))
			set[t] = true
		}
		nonuniquecorpus = append(nonuniquecorpus, strings.ToLower(t))
		x, _ := megophone.Metaphone(t)
		nonuniquecorpus = append(nonuniquecorpus, strings.ToLower(x))
		stem := porterstemmer.Stem([]rune(t))
		nonuniquecorpus = append(nonuniquecorpus, strings.ToLower(string(stem)))
	}
	sort.Strings(uniquecorpus)
	return uniquecorpus, nonuniquecorpus
}

func wordFreq(wordList []string) map[string]int {
	counts := make(map[string]int)
	for _, word := range wordList {
		_, ok := counts[word]
		if ok {
			counts[word]++
		} else {
			counts[word] = 1
		}
	}
	return counts
}

func cosineSimilarity(a []float64, b []float64) (cosine float64, err error) {
	count := 0
	lengtha := len(a)
	lengthb := len(b)
	if lengtha > lengthb {
		count = lengtha
	} else {
		count = lengthb
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= lengtha {
			s2 += math.Pow(b[k], 2)
			continue
		}
		if k >= lengthb {
			s1 += math.Pow(a[k], 2)
			continue
		}
		sumA += a[k] * b[k]
		s1 += math.Pow(a[k], 2)
		s2 += math.Pow(b[k], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, errors.New("Vectors should not be null (all zeros)")
	}
	return sumA / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}

func absoluteSimilarity(a []float64, b []float64) float64 {
	nonZeroIntersections := float64(0)
	for i := range a {
		if a[i] > 0 && b[i] > 0 {
			nonZeroIntersections++
		}
	}
	return nonZeroIntersections / float64(len(a))
}
