package ccprotos

func (con *Contract) Find(addr string) (*Contract_MemberStatus, bool) {

	for _, v := range con.Status {
		if v.MemberID == addr {
			return v, true
		}
	}

	return nil, false
}
