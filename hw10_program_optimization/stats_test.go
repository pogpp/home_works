//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
			  {"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
			  {"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
			  {"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
			  {"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	data2 := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@я.ру","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
			  {"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@я.ру","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`

	data3 := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
			  {"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
			  {"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}`

	data4 := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"@browsedrive.su","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
			  {"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"@browsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
			  {"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@browsecat.ru","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}`

	data5 := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
			  {"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@browsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
			  {"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("find empty domain", func(t *testing.T) {
		_, err := GetDomainStat(bytes.NewBufferString(data), "")
		require.ErrorIs(t, err, ErrEmptyDomain)
	})

	t.Run("find domain with capital letter", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "Com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find cyrillic 'ру'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data2), "ру")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"я.ру": 2,
		}, result)
	})

	t.Run("find 'com' without domain", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data3), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 1,
		}, result)
	})

	t.Run("find 'com' without user name", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data4), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 1,
		}, result)
	})

	t.Run("find 'com' none symbols after @", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data5), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 1,
		}, result)
	})
}
