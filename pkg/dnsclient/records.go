package dnsclient

import (
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
)

// added constant for DNS record types
const (
	type_A     = "A"
	type_CNAME = "CNAME"
	type_AAAA  = "AAAA"
	type_TXT   = "TXT"

	username = "username"
	password = "password"
)

type Record interface {
	GetId() string
	GetType() string
	GetValue() string
	GetDNSName() string
	GetSetIdentifier() string
	GetTTL() int
	SetTTL(int)
	Copy() Record
}

type RecordA ibclient.RecordA

func (r *RecordA) GetType() string          { return type_A }
func (r *RecordA) GetId() string            { return r.Ref }
func (r *RecordA) GetDNSName() string       { return r.Name }
func (r *RecordA) GetSetIdentifier() string { return "" }
func (r *RecordA) GetValue() string         { return r.Ipv4Addr }
func (r *RecordA) GetTTL() int              { return int(r.Ttl) }
func (r *RecordA) SetTTL(ttl int)           { r.Ttl = uint32(ttl); r.UseTtl = ttl != 0 }
func (r *RecordA) Copy() Record             { n := *r; return &n }
func (r *RecordA) PrepareUpdate() Record {
	n := *r
	n.Zone = ""
	n.Name = ""
	n.View = ""
	return &n
}

type RecordAAAA ibclient.RecordAAAA

func (r *RecordAAAA) GetType() string          { return type_AAAA }
func (r *RecordAAAA) GetId() string            { return r.Ref }
func (r *RecordAAAA) GetDNSName() string       { return r.Name }
func (r *RecordAAAA) GetSetIdentifier() string { return "" }
func (r *RecordAAAA) GetValue() string         { return r.Ipv6Addr }
func (r *RecordAAAA) GetTTL() int              { return int(r.Ttl) }
func (r *RecordAAAA) SetTTL(ttl int)           { r.Ttl = uint32(ttl); r.UseTtl = ttl != 0 }
func (r *RecordAAAA) Copy() Record             { n := *r; return &n }
func (r *RecordAAAA) PrepareUpdate() Record {
	n := *r
	n.Zone = ""
	n.Name = ""
	n.View = ""
	return &n
}

type RecordCNAME ibclient.RecordCNAME

func (r *RecordCNAME) GetType() string          { return type_CNAME }
func (r *RecordCNAME) GetId() string            { return r.Ref }
func (r *RecordCNAME) GetDNSName() string       { return r.Name }
func (r *RecordCNAME) GetSetIdentifier() string { return "" }
func (r *RecordCNAME) GetValue() string         { return r.Canonical }
func (r *RecordCNAME) GetTTL() int              { return int(r.Ttl) }
func (r *RecordCNAME) SetTTL(ttl int)           { r.Ttl = uint32(ttl); r.UseTtl = ttl != 0 }
func (r *RecordCNAME) Copy() Record             { n := *r; return &n }
func (r *RecordCNAME) PrepareUpdate() Record    { n := *r; n.Zone = ""; n.View = ""; return &n }

type RecordTXT ibclient.RecordTXT

func (r *RecordTXT) GetType() string          { return type_TXT }
func (r *RecordTXT) GetId() string            { return r.Ref }
func (r *RecordTXT) GetDNSName() string       { return r.Name }
func (r *RecordTXT) GetSetIdentifier() string { return "" }
func (r *RecordTXT) GetValue() string         { return EnsureQuotedText(r.Text) }
func (r *RecordTXT) GetTTL() int              { return int(r.Ttl) }
func (r *RecordTXT) SetTTL(ttl int)           { r.Ttl = uint(ttl); r.UseTtl = ttl != 0 }
func (r *RecordTXT) Copy() Record             { n := *r; return &n }
func (r *RecordTXT) PrepareUpdate() Record    { n := *r; n.Zone = ""; n.View = ""; return &n }

var _ Record = (*RecordA)(nil)
var _ Record = (*RecordAAAA)(nil)
var _ Record = (*RecordCNAME)(nil)
var _ Record = (*RecordTXT)(nil)

type RecordNS ibclient.RecordNS
