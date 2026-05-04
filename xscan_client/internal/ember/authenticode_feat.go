package ember

import (
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"time"
	"unicode/utf16"

	"github.com/ayoubfaouzi/pkcs7"
	"github.com/saferwall/pe"
)

// OIDs aligned with signify / asn1crypto cms for Authenticode (ember_cert features.py).
var (
	oidSpcSpOpusInfo    = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 1, 12}
	oidNestedSignature  = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 2, 4, 1}
	oidCounterSignature = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 6}
	oidSigningTime      = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 5}
	oidMSCertTimestamp  = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 3, 3, 1}
)

type cmsAttr struct {
	Type  asn1.ObjectIdentifier
	Value asn1.RawValue `asn1:"set"`
}

type signerDERWalk struct {
	Version                   int
	IssuerAndSerialNumber     asn1.RawValue
	DigestAlgorithm           asn1.RawValue
	AuthenticatedAttributes   asn1.RawValue `asn1:"optional,implicit,tag:0"`
	DigestEncryptionAlgorithm asn1.RawValue
	EncryptedDigest           []byte
	UnauthenticatedAttributes asn1.RawValue `asn1:"optional,implicit,tag:1"`
}

type signedDataOuter struct {
	Version                    int
	DigestAlgorithmIdentifiers []asn1.RawValue `asn1:"set"`
	ContentInfo                asn1.RawValue
	Certificates               asn1.RawValue   `asn1:"optional,tag:0"`
	CRLs                       asn1.RawValue   `asn1:"optional,tag:1"`
	SignerInfos                []signerDERWalk `asn1:"set"`
}

type contentInfoOuter struct {
	ContentType asn1.ObjectIdentifier
	Content     asn1.RawValue `asn1:"explicit,optional,tag:0"`
}

type authAccum struct {
	numSigs    int
	selfSigned float32
	emptyProg  float32
	noCounter  float32
	parseErr   float32
	chainMax   int
	latestSig  float64
	sigDiff    float64
}

// ProcessAuthenticodeSignature matches go_server/ember_cert/features.py AuthenticodeSignature
// (signify SignedPEFile.iter_embedded_signatures + raw_features / process_raw_features).
func ProcessAuthenticodeSignature(p *pe.File, includeCert bool, dst []float32) {
	if len(dst) < DimAuthenticode {
		return
	}
	clear(dst[:DimAuthenticode])
	if !includeCert || p == nil || len(p.Certificates.Raw) == 0 {
		return
	}

	peStamp := p.NtHeader.FileHeader.TimeDateStamp
	var acc authAccum
	if err := walkAuthenticodeRecursive(p.Certificates.Raw, peStamp, &acc); err != nil {
		dst[4] = 1 // parse_error
		return
	}

	dst[0] = float32(acc.numSigs)
	dst[1] = acc.selfSigned
	dst[2] = acc.emptyProg
	dst[3] = acc.noCounter
	dst[4] = acc.parseErr
	dst[5] = float32(acc.chainMax)
	dst[6] = float32(acc.latestSig)
	dst[7] = float32(acc.sigDiff)
}

func walkAuthenticodeRecursive(raw []byte, peStamp uint32, acc *authAccum) error {
	p7, err := pkcs7.Parse(raw)
	if err != nil {
		return err
	}

	acc.numSigs++
	updateCertChainFeatures(p7, acc)

	if programNameMissingLikeSignify(p7) {
		acc.emptyProg = 1
	}

	ts, ok := extractCountersignerUnix(p7)
	if !ok {
		acc.noCounter = 1
	} else {
		if ts >= acc.latestSig {
			acc.latestSig = ts
		}
		acc.sigDiff = ts - float64(peStamp)
	}

	nestedDERs := nestedAuthenticodeDERList(raw)
	for _, nd := range nestedDERs {
		_ = walkAuthenticodeRecursive(nd, peStamp, acc)
	}
	return nil
}

func updateCertChainFeatures(p7 *pkcs7.PKCS7, acc *authAccum) {
	n := len(p7.Certificates)
	if n > acc.chainMax {
		acc.chainMax = n
	}
	for i := 0; i < n-1; i++ {
		c := p7.Certificates[i]
		if c.Issuer.String() == c.Subject.String() {
			acc.selfSigned = 1
		}
	}
}

func programNameMissingLikeSignify(p7 *pkcs7.PKCS7) bool {
	var raw asn1.RawValue
	err := p7.UnmarshalSignedAttribute(oidSpcSpOpusInfo, &raw)
	if err != nil {
		var nf pkcs7.AttributeNotFoundError
		if errors.As(err, &nf) {
			return true
		}
		return true
	}

	var opus struct {
		ProgramName *asn1.RawValue `asn1:"optional,explicit,tag:0"`
	}
	if _, err := asn1.Unmarshal(raw.Bytes, &opus); err != nil {
		return true
	}
	if opus.ProgramName == nil {
		return true
	}
	s, err := decodeSpcString(*opus.ProgramName)
	if err != nil || s == "" {
		return true
	}
	return false
}

func decodeSpcString(explicit asn1.RawValue) (string, error) {
	var choice asn1.RawValue
	if _, err := asn1.Unmarshal(explicit.Bytes, &choice); err != nil {
		return "", err
	}
	switch {
	case choice.Class == 2 && choice.Tag == 0:
		return utf16BEToString(choice.Bytes), nil
	case choice.Class == 2 && choice.Tag == 1:
		return string(choice.Bytes), nil
	default:
		var inner asn1.RawValue
		if _, err := asn1.Unmarshal(explicit.Bytes, &inner); err != nil {
			return "", err
		}
		if inner.Tag == 30 {
			return utf16BEToString(inner.Bytes), nil
		}
		return string(inner.Bytes), nil
	}
}

func utf16BEToString(b []byte) string {
	if len(b) < 2 {
		return ""
	}
	n := len(b) / 2
	u := make([]uint16, n)
	for i := 0; i < n; i++ {
		u[i] = binary.BigEndian.Uint16(b[i*2:])
	}
	return string(utf16.Decode(u))
}

func extractCountersignerUnix(p7 *pkcs7.PKCS7) (float64, bool) {
	var csRaw asn1.RawValue
	if err := p7.UnmarshalUnsignedAttribute(oidCounterSignature, &csRaw); err == nil {
		if t, err := signingTimeFromPKCS9CounterSignature(csRaw); err == nil {
			return unixFloat(t), true
		}
	}
	var tsRaw asn1.RawValue
	if err := p7.UnmarshalUnsignedAttribute(oidMSCertTimestamp, &tsRaw); err == nil {
		if t, err := signingTimeFromMSRFC3161(tsRaw); err == nil {
			return unixFloat(t), true
		}
	}
	return 0, false
}

func unixFloat(t time.Time) float64 {
	return float64(t.Unix())
}

func signingTimeFromPKCS9CounterSignature(attrVal asn1.RawValue) (time.Time, error) {
	var sw signerDERWalk
	var err error
	switch attrVal.Tag {
	case 17: // SET OF SignerInfo
		var inner asn1.RawValue
		if _, err = asn1.Unmarshal(attrVal.Bytes, &inner); err != nil {
			return time.Time{}, err
		}
		_, err = asn1.Unmarshal(inner.FullBytes, &sw)
	case 16: // some encoders use a single SignerInfo SEQUENCE (matches Chrom.exe)
		_, err = asn1.Unmarshal(attrVal.FullBytes, &sw)
	default:
		_, err = asn1.Unmarshal(attrVal.FullBytes, &sw)
	}
	if err != nil {
		return time.Time{}, err
	}
	attrs := parseCMSAttributes(sw.AuthenticatedAttributes)
	for _, a := range attrs {
		if a.Type.Equal(oidSigningTime) {
			return parsePKCS9SigningTimeValue(a.Value)
		}
	}
	return time.Time{}, errors.New("no signing_time in countersigner")
}

func signingTimeFromMSRFC3161(attrVal asn1.RawValue) (time.Time, error) {
	var set asn1.RawValue
	if _, err := asn1.Unmarshal(attrVal.Bytes, &set); err != nil {
		return time.Time{}, err
	}
	var ciDER asn1.RawValue
	if _, err := asn1.Unmarshal(set.Bytes, &ciDER); err != nil {
		return time.Time{}, err
	}
	p7, err := pkcs7.Parse(ciDER.FullBytes)
	if err != nil {
		return time.Time{}, err
	}
	var tst struct {
		Version        int
		Policy         asn1.ObjectIdentifier
		MessageImprint asn1.RawValue
		SerialNumber   asn1.RawValue
		GenTime        time.Time `asn1:"generalized"`
	}
	if _, err := asn1.Unmarshal(p7.Content, &tst); err != nil {
		return time.Time{}, err
	}
	return tst.GenTime, nil
}

func parsePKCS9SigningTimeValue(val asn1.RawValue) (time.Time, error) {
	var inner asn1.RawValue
	if _, err := asn1.Unmarshal(val.Bytes, &inner); err != nil {
		return time.Time{}, err
	}
	switch inner.Tag {
	case 23, 24:
		var t time.Time
		if _, err := asn1.Unmarshal(inner.FullBytes, &t); err != nil {
			return time.Time{}, err
		}
		return t, nil
	default:
		var t time.Time
		if _, err := asn1.Unmarshal(inner.Bytes, &t); err != nil {
			return time.Time{}, err
		}
		return t, nil
	}
}

func parseCMSAttributes(block asn1.RawValue) []cmsAttr {
	if len(block.Bytes) == 0 && len(block.FullBytes) == 0 {
		return nil
	}
	data := block.Bytes
	var set asn1.RawValue
	if _, err := asn1.Unmarshal(block.FullBytes, &set); err == nil && set.Tag == 17 {
		data = set.Bytes
	} else if _, err := asn1.Unmarshal(data, &set); err == nil && set.Tag == 17 {
		data = set.Bytes
	}
	var attrs []cmsAttr
	for len(data) > 0 {
		var a cmsAttr
		rest, err := asn1.Unmarshal(data, &a)
		if err != nil {
			break
		}
		attrs = append(attrs, a)
		data = rest
	}
	return attrs
}

func nestedAuthenticodeDERList(raw []byte) [][]byte {
	var ci contentInfoOuter
	if _, err := asn1.Unmarshal(raw, &ci); err != nil {
		return nil
	}
	var sd signedDataOuter
	if _, err := asn1.Unmarshal(ci.Content.Bytes, &sd); err != nil {
		return nil
	}
	if len(sd.SignerInfos) == 0 {
		return nil
	}
	unsigned := parseCMSAttributes(sd.SignerInfos[0].UnauthenticatedAttributes)
	var out [][]byte
	for _, a := range unsigned {
		if !a.Type.Equal(oidNestedSignature) {
			continue
		}
		parts := unpackSETOfDER(a.Value.Bytes)
		out = append(out, parts...)
	}
	return out
}

func unpackSETOfDER(setBytes []byte) [][]byte {
	var outer asn1.RawValue
	if _, err := asn1.Unmarshal(setBytes, &outer); err != nil {
		return nil
	}
	inner := outer.Bytes
	var seqs [][]byte
	for len(inner) > 0 {
		var item asn1.RawValue
		rest, err := asn1.Unmarshal(inner, &item)
		if err != nil {
			break
		}
		seqs = append(seqs, append([]byte(nil), item.FullBytes...))
		inner = rest
	}
	return seqs
}
