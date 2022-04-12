package globals

import (
	types "blockchain/blockchaintypes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"time"
)

type GlobalLib struct{}

type IGlobalLib interface {
	IsHttpOk(err error, w http.ResponseWriter) bool
	NowUnixNano() int64
	EmptyByte32() types.Byte32
	JsonStatus(message string) []byte
	DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error
	PublicKeyFromString(s string) *ecdsa.PublicKey
	PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey
	SignatureFromString(s string) *Signature
}

func NewGlobals() IGlobalLib {
	return &GlobalLib{}
}

func (g *GlobalLib) NowUnixNano() int64 {
	return time.Now().UnixNano()
}

func (g *GlobalLib) IsHttpOk(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = io.WriteString(w, string(g.JsonStatus(fmt.Sprintf("failure: %v", err))))
		return false
	} else {
		return true
	}
}

func (g *GlobalLib) String2BigIntTuple(s string) (big.Int, big.Int) {
	hx, _ := hex.DecodeString(s[:64])
	hy, _ := hex.DecodeString(s[64:])

	var bix big.Int
	var biy big.Int

	_ = bix.SetBytes(hx)
	_ = biy.SetBytes(hy)

	return bix, biy
}

func (g *GlobalLib) PublicKeyFromString(s string) *ecdsa.PublicKey {
	x, y := g.String2BigIntTuple(s)
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &x,
		Y:     &y,
	}
}

func (g *GlobalLib) PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(s[:])
	var bi big.Int
	_ = bi.SetBytes(b)
	return &ecdsa.PrivateKey{
		PublicKey: *publicKey,
		D:         &bi,
	}
}

func (g *GlobalLib) SignatureFromString(s string) *Signature {
	x, y := g.String2BigIntTuple(s)
	return &Signature{
		R: &x,
		S: &y,
	}
}

func (g *GlobalLib) EmptyByte32() types.Byte32 {
	return types.Byte32{}
}

func (g *GlobalLib) JsonStatus(message string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return m
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

// DecodeJSONBody
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func (g *GlobalLib) DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

//func UnmarshalRequestBody[T any](w http.ResponseWriter, req *http.Request) (T, error) {
//	// read the body
//	b, err := ioutil.ReadAll(req.Body)
//
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Println(err)
//		}
//	}(req.Body)
//
//	var t T
//
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//		return t, err
//	}
//
//	// Unmarshal
//	err = json.Unmarshal(b, &t)
//	if err != nil {
//		http.Error(w, err.Error(), 500)
//		return t, err
//	}
//
//	return t, nil
//}
//
//func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
//	if r.Header.Get("Content-Type") != "" {
//		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
//		if value != "application/json" {
//			msg := "Content-Type header is not application/json"
//			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
//		}
//	}
//
//	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
//
//	dec := json.NewDecoder(r.Body)
//	dec.DisallowUnknownFields()
//
//	err := dec.Decode(&dst)
//	if err != nil {
//		var syntaxError *json.SyntaxError
//		var unmarshalTypeError *json.UnmarshalTypeError
//
//		switch {
//		case errors.As(err, &syntaxError):
//			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
//			return &malformedRequest{status: http.StatusBadRequest, msg: msg}
//
//		case errors.Is(err, io.ErrUnexpectedEOF):
//			msg := fmt.Sprintf("Request body contains badly-formed JSON")
//			return &malformedRequest{status: http.StatusBadRequest, msg: msg}
//
//		case errors.As(err, &unmarshalTypeError):
//			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
//			return &malformedRequest{status: http.StatusBadRequest, msg: msg}
//
//		case strings.HasPrefix(err.Error(), "json: unknown field "):
//			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
//			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
//			return &malformedRequest{status: http.StatusBadRequest, msg: msg}
//
//		case errors.Is(err, io.EOF):
//			msg := "Request body must not be empty"
//			return &malformedRequest{status: http.StatusBadRequest, msg: msg}
//
//		case err.Error() == "http: request body too large":
//			msg := "Request body must not be larger than 1MB"
//			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}
//
//		default:
//			return err
//		}
//	}
//
//	err = dec.Decode(&struct{}{})
//	if err != io.EOF {
//		msg := "Request body must only contain a single JSON object"
//		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
//	}
//
//	return nil
//}
