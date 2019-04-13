package helpers

import (
	"strconv"
	"strings"

	"github.com/OneOfOne/xxhash"
)

type addressType string

func (t addressType) ToString() string {
	return string(t)
}

const (
	incomeAddressType addressType = "income"
	spendAddressType addressType = "spend"
)

const (
	commonAddressPrefix string = "0x0b"
	incomeAddressPrefix string = "0x0d"
	spendAddressPrefix string = "0x0c"
)

const (
	addressHashSeed		= 0x2018CADD
	addressHashBase		= 16
	oldAddressLen		= 44
)

func createAddressFromString(source string) string {
	return createAddressFromBytes([]byte(source))
}

func createAddressFromBytes(source []byte) string {
	sha := string(CreateHash(source)[:])

	return sha
}

func createSpecialAddress(trackerPublicKey string, peerPublicKey string, addressType addressType, timestamp int64) string {
	return createAddressFromString(trackerPublicKey + ";" + peerPublicKey +";" + addressType.ToString() + ";" + strconv.FormatInt(timestamp, 16))
}

func createAddressHash(prefix string, addressValue string) string {
	return strconv.FormatUint(xxhash.ChecksumString64S(prefix + addressValue, addressHashSeed), addressHashBase)
}

func CreateCommonAddress(pk string) string {
	addr := createAddressFromString(pk)
	hash := createAddressHash(commonAddressPrefix, addr)
	return commonAddressPrefix + hash + addr
}

func CreateIncomeAddress(trackerPublicKey string, peerPublicKey string, timestamp int64) string {
	addr := createSpecialAddress(trackerPublicKey, peerPublicKey, incomeAddressType, timestamp)
	hash := createAddressHash(incomeAddressPrefix, addr)
	return incomeAddressPrefix + hash + addr
}

func CreateSpendAddress(trackerPublicKey string, peerPublicKey string, timestamp int64) string {
	addr := createSpecialAddress(trackerPublicKey, peerPublicKey, spendAddressType, timestamp)
	hash := createAddressHash(spendAddressPrefix, addr)
	return spendAddressPrefix + hash + addr
}

func extractHashAddress(prefix string, runes []rune) (address string, hash string) {
	prefixLen := len([]rune(prefix))
	addressLen := len(runes)
	pos := addressLen - oldAddressLen + prefixLen
	hash = string(runes[prefixLen:pos])
	address = string(runes[pos:addressLen])
	return
}

func isAddressValid(prefix string, address string) bool {
	isRightPrefix := strings.HasPrefix(address, prefix)
	if !isRightPrefix {
		return false
	}

	// Condition for backward compability
	runes := []rune(address)
	if len(runes) == oldAddressLen {
		return isRightPrefix
	}

	addr, hash := extractHashAddress(prefix, runes)
	return hash == createAddressHash(prefix, addr)

}

func IsAddressCommon(address string) bool {
	return isAddressValid(commonAddressPrefix, address)
}

func IsAddressSpend(address string) bool {
	return isAddressValid(spendAddressPrefix, address)
}

func IsAddressIncome(address string) bool {
	return isAddressValid(incomeAddressPrefix, address)
}
