// jwt is a command line utility that encodes or decodes a JSON Web Token,
// using provided key data. jwt can use any PEM encoded public or private key
// file, including JSON-encoded keys (such as a Google service account
// credential file) to encode, decode, and verify JWTs.
//
// Example:
//    # encode arbitrary JSON as payload (ie, claims)
//    echo '{"iss": "issuer", "nbf": '$(date +%s)'}' | jwt -k ./testdata/rsa.pem -enc
//
//    # encode name/value pairs from command line
//    jwt -k ./testdata/rsa.pem -enc iss=issuer nbf=$(date +%s)
//
//    # decode (and verify) token
//    echo "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg" | jwt -k ./testdata/rsa.pem -dec
//
//    # encode and decode in one sweep:
//    jwt -k ./testdata/rsa.pem -enc iss=issuer nbf=$(date +%s) | jwt -k ./testdata/rsa.pem -dec
//
//    # specify algorithm -- this will error since the token here is encoded using RS256, not RS384
//    echo "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg" | jwt -k ./testdata/rsa.pem -dec -alg RS384

package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"

	"github.com/knq/jwt"
	"github.com/knq/pemutil"
)

var (
	flagEnc = flag.Bool("enc", false, "encode token from json data provided from stdin, or via name=value pairs passed on the command line")
	flagDec = flag.Bool("dec", false, "decode and verify token read from stdin using the provided key data")
	flagKey = flag.String("k", "", "path to PEM-encoded file or JSON file containing key data")
	flagAlg = flag.String("alg", "", "override signing algorithm")
)

func main() {
	var err error

	// parse parameters
	flag.Parse()

	// make sure k parameter is specified
	if *flagKey == "" {
		fmt.Fprintln(os.Stderr, "error: must supply a key file with -k")
		os.Exit(1)
	}

	// inspect remaining args
	args := flag.Args()
	if len(args) > 0 && *flagDec {
		fmt.Fprintln(os.Stderr, "error: unknown args passed for -dec")
		os.Exit(1)
	}

	// if there are command line args and enc, then build js from them
	var in []byte
	if len(args) > 0 && *flagEnc {
		in, err = buildEncArgs(args)
	} else {
		in, err = ioutil.ReadAll(os.Stdin)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// read key data and suggestedAlg
	pem, suggestedAlg, err := loadKeyData(*flagKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// determine alg
	var alg jwt.Algorithm
	if suggestedAlg != jwt.NONE {
		alg = suggestedAlg
	} else if *flagAlg != "" {
		err = (&alg).UnmarshalText([]byte(*flagAlg))
	} else if *flagDec {
		alg, err = jwt.PeekAlgorithm(in)
	} else {
		alg, err = getAlgFromKeyData(pem)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// create signer
	signer, err := alg.New(pem)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// encode or decode
	var out []byte
	switch {
	case *flagDec:
		out, err = doDec(signer, in)

	case *flagEnc:
		out, err = doEnc(signer, in)

	default:
		err = errors.New("please specify -enc or -dec")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// make the output a little nicer
	if isatty.IsTerminal(os.Stdout.Fd()) {
		out = append(out, '\n')
	}

	os.Stdout.Write(out)
}

// loadKeyData loads key data from the path.
//
// Attempts to json-decode the file first, otherwise passes it to pemutil. If
// the file appears to be a Google service account JSON file (ie contains
// "gserviceaccount.com"), then jwt.RS256 will be returned, otherwise jwt.NONE
// is returned.
func loadKeyData(path string) (pemutil.Store, jwt.Algorithm, error) {
	// check if it's json data
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, jwt.NONE, err
	}

	// attempt json decode
	v := make(map[string]interface{})
	err = json.Unmarshal(buf, &v)
	if err != nil {
		// not json encoded, so skip
		pem := pemutil.Store{}
		err = pemutil.PEM{path}.Load(pem)
		if err != nil {
			return nil, jwt.NONE, err
		}
		return pem, jwt.NONE, nil
	}

	// attempt decode on each field of the json decoded values, and ignoring
	// any errors
	pem := pemutil.Store{}
	for _, val := range v {
		if str, ok := val.(string); ok {
			_ = pemutil.DecodePEM(pem, []byte(str))
		}
	}

	if len(pem) < 1 {
		return nil, jwt.NONE, fmt.Errorf("could not load any PEM encoded keys from %s", path)
	}

	// if it's a google service account, return RS256
	if bytes.Contains(buf, []byte("gserviceaccount.com")) {
		return pem, jwt.RS256, nil
	}

	return pem, jwt.NONE, nil
}

// getSuitableAlgFromCurve inspects the key length in curve, and determines the
// corresponding jwt.Algorithm.
func getSuitableAlgFromCurve(curve elliptic.Curve) (jwt.Algorithm, error) {
	curveBitSize := curve.Params().BitSize

	// compute curve key len
	keyLen := curveBitSize / 8
	if curveBitSize%8 > 0 {
		keyLen++
	}

	// determine alg
	var alg jwt.Algorithm
	switch 2 * keyLen {
	case 64:
		alg = jwt.ES256
	case 96:
		alg = jwt.ES384
	case 132:
		alg = jwt.ES512

	default:
		return jwt.NONE, fmt.Errorf("invalid key length %d", keyLen)
	}

	return alg, nil
}

// getAlgFromKeyData determines the best jwt.Algorithm suitable based on the
// set of given crypto primitives in pem.
func getAlgFromKeyData(pem pemutil.Store) (jwt.Algorithm, error) {
	for _, v := range pem {
		// loop over crypto primitives in pemstore, and do type assertion. if
		// ecdsa.{PublicKey,PrivateKey} found, then use corresponding ESXXX as
		// algo. if rsa, then use DefaultRSAAlgorithm. if []byte, then use
		// DefaultHMACAlgorithm.
		switch k := v.(type) {
		case []byte:
			return jwt.HS512, nil

		case *ecdsa.PrivateKey:
			return getSuitableAlgFromCurve(k.Curve)

		case *ecdsa.PublicKey:
			return getSuitableAlgFromCurve(k.Curve)

		case *rsa.PrivateKey:
			return jwt.PS512, nil

		case *rsa.PublicKey:
			return jwt.PS512, nil
		}
	}

	return jwt.NONE, errors.New("cannot determine key type")
}

// buildEncArgs builds and encodes passed argument strings in the form of
// name=val as a json object.
func buildEncArgs(args []string) ([]byte, error) {
	m := make(map[string]interface{})

	// loop over args, splitting on '=', and attempt parsing of value
	for _, arg := range args {
		a := strings.SplitN(arg, "=", 2)
		var val interface{}

		// attempt to parse
		if len(a) == 1 { // assume bool, set as true
			val = true
		} else if u, err := strconv.ParseUint(a[1], 10, 64); err == nil {
			val = u
		} else if i, err := strconv.ParseInt(a[1], 10, 64); err == nil {
			val = i
		} else if f, err := strconv.ParseFloat(a[1], 64); err == nil {
			val = f
		} else if b, err := strconv.ParseBool(a[1]); err == nil {
			val = b
		} else if s, err := strconv.Unquote(a[1]); err == nil {
			val = s
		} else { // treat as string
			val = a[1]
		}

		m[a[0]] = val
	}

	return json.Marshal(m)
}

// UnstructuredToken is a jwt compatible token for encoding/decoding unknown
// jwt payloads.
type UnstructuredToken struct {
	Header    map[string]interface{} `json:"header" jwt:"header"`
	Payload   map[string]interface{} `json:"payload" jwt:"payload"`
	Signature []byte                 `json:"signature" jwt:"signature"`
}

// doDec decodes in as a JWT.
func doDec(signer jwt.Signer, in []byte) ([]byte, error) {
	var err error

	// decode token
	ut := UnstructuredToken{}
	err = signer.Decode(bytes.TrimSpace(in), &ut)
	if err != nil {
		return nil, err
	}

	// pretty format output
	out, err := json.MarshalIndent(&ut, "", "  ")
	if err != nil {
		return nil, err
	}

	return out, nil
}

// doEnc encodes in as the payload in a JWT.
func doEnc(signer jwt.Signer, in []byte) ([]byte, error) {
	var err error

	// make sure its valid json first
	m := make(map[string]interface{})

	// do the initial decode
	d := json.NewDecoder(bytes.NewBuffer(in))
	d.UseNumber()
	err = d.Decode(&m)
	if err != nil {
		return nil, err
	}

	// encode claims
	out, err := signer.Encode(&m)
	if err != nil {
		return nil, err
	}

	return out, nil
}
