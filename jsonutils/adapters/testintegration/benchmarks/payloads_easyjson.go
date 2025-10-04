// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package benchmarks

import (
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

const defaultEasyJSONAlloc = 8

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SmallPayload) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	first := true
	{
		const prefix string = ",\"st\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(v.St)
	}
	{
		const prefix string = ",\"sid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(v.Sid)
	}
	{
		const prefix string = ",\"tt\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(v.Tt)
	}
	{
		const prefix string = ",\"gr\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(v.Gr)
	}
	{
		const prefix string = ",\"uuid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(v.UUID)
	}
	{
		const prefix string = ",\"ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(v.IP)
	}
	{
		const prefix string = ",\"ua\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(v.Ua)
	}
	{
		const prefix string = ",\"tz\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(v.Tz)
	}
	{
		const prefix string = ",\"v\":"
		if first {
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(v.V)
	}
	out.RawByte('}')
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SmallPayload) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "st":
			v.St = in.Int()
		case "sid":
			v.Sid = in.Int()
		case "tt":
			v.Tt = in.String()
		case "gr":
			v.Gr = in.Int()
		case "uuid":
			v.UUID = in.String()
		case "ip":
			v.IP = in.String()
		case "ua":
			v.Ua = in.String()
		case "tz":
			v.Tz = in.Int()
		case "v":
			v.V = in.Int()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v MediumPayload) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	first := true
	{
		const prefix string = ",\"person\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if v.Person == nil {
			out.RawString("null")
		} else {
			v.Person.MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"company\":"
		if first {
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(v.Company)
	}
	out.RawByte('}')
}

func (v *MediumPayload) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "person":
			if in.IsNull() {
				in.Skip()
				v.Person = nil
			} else {
				if v.Person == nil {
					v.Person = new(CBPerson)
				}
				v.Person.UnmarshalEasyJSON(in)
			}
		case "company":
			v.Company = in.String()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v CBPerson) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	first := true
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if v.Name == nil {
			out.RawString("null")
		} else {
			v.Name.MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"github\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if v.Github == nil {
			out.RawString("null")
		} else {
			v.Github.MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"gravatar\":"
		if first {
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if v.Gravatar == nil {
			out.RawString("null")
		} else {
			v.Gravatar.MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

func (v *CBPerson) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			if in.IsNull() {
				in.Skip()
				v.Name = nil
			} else {
				if v.Name == nil {
					v.Name = new(CBName)
				}
				v.Name.UnmarshalEasyJSON(in)
			}
		case "github":
			if in.IsNull() {
				in.Skip()
				v.Github = nil
			} else {
				if v.Github == nil {
					v.Github = new(CBGithub)
				}
				v.Github.UnmarshalEasyJSON(in)
			}
		case "gravatar":
			if in.IsNull() {
				in.Skip()
				v.Gravatar = nil
			} else {
				if v.Gravatar == nil {
					v.Gravatar = new(CBGravatar)
				}
				v.Gravatar.UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v CBName) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	first := true
	{
		const prefix string = ",\"fullName\":"
		if first {
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(v.FullName)
	}
	out.RawByte('}')
}

func (v *CBName) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "fullName":
			v.FullName = in.String()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v CBGithub) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	{
		const prefix string = ",\"followers\":"
		out.RawString(prefix[1:])
		out.Int(v.Followers)
	}
	out.RawByte('}')
}

func (v *CBGithub) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "followers":
			v.Followers = in.Int()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v CBGravatar) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	{
		const prefix string = ",\"avatars\":"
		out.RawString(prefix[1:])

		if v.Avatars == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range v.Avatars {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					v3.MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

func (v *CBGravatar) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "avatars":
			if in.IsNull() {
				in.Skip()
				v.Avatars = nil
			} else {
				in.Delim('[')
				if v.Avatars == nil {
					if !in.IsDelim(']') {
						v.Avatars = make(Avatars, 0, defaultEasyJSONAlloc)
					} else {
						v.Avatars = Avatars{}
					}
				} else {
					v.Avatars = (v.Avatars)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *CBAvatar
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(CBAvatar)
						}
						v1.UnmarshalEasyJSON(in)
					}
					v.Avatars = append(v.Avatars, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v CBAvatar) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	{
		const prefix string = ",\"url\":"
		out.RawString(prefix[1:])
		out.String(v.URL)
	}
	out.RawByte('}')
}

func (v *CBAvatar) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "url":
			v.URL = in.String()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v LargePayload) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	first := true
	{
		const prefix string = ",\"users\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if v.Users == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range v.Users {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					v3.MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"topics\":"
		if first {
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if v.Topics == nil {
			out.RawString("null")
		} else {
			v.Topics.MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

func (v *LargePayload) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "users":
			if in.IsNull() {
				in.Skip()
				v.Users = nil
			} else {
				in.Delim('[')
				if v.Users == nil {
					if !in.IsDelim(']') {
						v.Users = make(DSUsers, 0, defaultEasyJSONAlloc)
					} else {
						v.Users = DSUsers{}
					}
				} else {
					v.Users = (v.Users)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *DSUser
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(DSUser)
						}
						v1.UnmarshalEasyJSON(in)
					}
					v.Users = append(v.Users, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "topics":
			if in.IsNull() {
				in.Skip()
				v.Topics = nil
			} else {
				if v.Topics == nil {
					v.Topics = new(DSTopicsList)
				}
				v.Topics.UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v DSUser) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix[1:])
		out.String(v.Username)
	}
	out.RawByte('}')
}

func (v *DSUser) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "username":
			v.Username = in.String()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v DSTopicsList) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	first := true
	{
		const prefix string = ",\"topics\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if v.Topics == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range v.Topics {
				if v5 > 0 {
					out.RawByte(',')
				}
				if v6 == nil {
					out.RawString("null")
				} else {
					v6.MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"moreTopicsURL\":"
		if first {
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(v.MoreTopicsURL)
	}
	out.RawByte('}')
}

func (v *DSTopicsList) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "topics":
			if in.IsNull() {
				in.Skip()
				v.Topics = nil
			} else {
				in.Delim('[')
				if v.Topics == nil {
					if !in.IsDelim(']') {
						v.Topics = make(DSTopics, 0, defaultEasyJSONAlloc)
					} else {
						v.Topics = DSTopics{}
					}
				} else {
					v.Topics = (v.Topics)[:0]
				}
				for !in.IsDelim(']') {
					var v4 *DSTopic
					if in.IsNull() {
						in.Skip()
						v4 = nil
					} else {
						if v4 == nil {
							v4 = new(DSTopic)
						}
						v4.UnmarshalEasyJSON(in)
					}
					v.Topics = append(v.Topics, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "moreTopicsURL":
			v.MoreTopicsURL = in.String()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func (v DSTopic) MarshalEasyJSON(out *jwriter.Writer) {
	out.RawByte('{')
	first := true
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(v.ID)
	}
	{
		const prefix string = ",\"slug\":"
		if first {
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(v.Slug)
	}
	out.RawByte('}')
}

func (v *DSTopic) UnmarshalEasyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			v.ID = in.Int()
		case "slug":
			v.Slug = in.String()
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
